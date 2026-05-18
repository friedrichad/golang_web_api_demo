package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// BasicAuthenticator validates client credentials (client_id and client_secret)
// Used for OAuth2 authentication endpoints
func BasicAuthenticator() gin.HandlerFunc {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(viper.GetString("oauth.client-id") + ":" + viper.GetString("oauth.client-secret")))
	return func(c *gin.Context) {
		basic := c.Request.Header["Authorization"]
		if len(basic) == 0 || strings.Replace(basic[0], "Basic ", "", 1) != basicAuth {
			c.JSON(http.StatusUnauthorized, model.ResponseWrapper{
				Code:    "401",
				Message: "Client không hợp lệ",
			})
			c.Abort()
			return
		}
	}
}

// BearerAuthenticator validates JWT bearer token
// Extracts user information and sets it in context for downstream handlers
func BearerAuthenticator() gin.HandlerFunc {
	hmacSecret := []byte(viper.GetString("oauth.jwt-secret"))
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Extract raw token without "Bearer " prefix for blacklist check
		rawToken := strings.Replace(authHeader, "Bearer ", "", 1)

		// Check if token is blacklisted (logout)
		isBlacklisted, err := redis.IsBlacklisted(rawToken)
		if err != nil {
			log.Printf("Error checking blacklist: %v", err)
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, model.ResponseWrapper{
				Code:    common.TokenInvalid.Code,
				Message: "Token đã bị đăng xuất",
			})
			c.Abort()
			return
		}

		jwtClaims, valid := extractClaims(authHeader, hmacSecret)
		if !valid {
			c.JSON(http.StatusUnauthorized, model.ResponseWrapper{
				Code:    common.TokenInvalid.Code,
				Message: common.TokenInvalid.Message,
			})
			c.Abort()
			return
		}
		c.Set("user_id", jwtClaims.Id)
		c.Set("username", jwtClaims.Username)
		c.Set("authorities", jwtClaims.Authorities)
		c.Set("position_level", jwtClaims.Level)
		c.Set("is_op", jwtClaims.IsOP)
		c.Next()
	}
}

// Authorizator checks if user has required authorities/permissions
// Check flow: JWT authorities -> Redis cache -> Database (lazy load + cache)
// Can be called with multiple authorities - user needs at least one of them
func Authorizator(authority ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetInt("is_op") == 1 {
			c.Next()
			return
		}
		jwtAuthorities := c.GetStringSlice("authorities")
		if len(jwtAuthorities) > 0 && utils.AnyContains(jwtAuthorities, authority) {
			c.Next()
			return
		}
		userIdStr := GetUserID(c)
		userId, _ := strconv.Atoi(userIdStr)

		// Check Redis cache first
		if userId > 0 && redis.CheckPermissionRedis(redis.Rdb, userId, authority) {
			c.Next()
			return
		}
		// If not in Redis, try to load from database and cache
		if userId > 0 {
			userRepo := repository.NewUserRepository()
			perms, err := userRepo.GetUserPermissionScopes(userId)
			if err != nil {
				log.Printf("Error fetching permissions from DB for user %d: %v", userId, err)
			} else if len(perms) > 0 {
				// Build permission map for Redis cache (scope -> expiredDate)
				permMap := make(map[string]interface{})
				for _, perm := range perms {
					permMap[perm.Scope] = perm.ExpiredDate
				}
				// Save to Redis cache (24 hour TTL)
				cacheErr := redis.SaveUserPermissionCache(redis.Rdb, userId, permMap, 24*time.Hour)
				if cacheErr != nil {
					log.Printf("Error saving permissions to Redis cache for user %d: %v", userId, cacheErr)
				}
				// Now check if requested authority exists in cached permissions
				if redis.CheckPermissionRedis(redis.Rdb, userId, authority) {
					c.Next()
					return
				}
			}
		}
		c.JSON(http.StatusForbidden, model.ResponseWrapper{
			Code:    "403",
			Message: "Không có quyền truy cập",
		})
		c.Abort()
	}
}

// GetUserID extracts user ID from context (set by BearerAuthenticator)
func GetUserID(c *gin.Context) string {
	userId, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userId.(string)
}

// GetUsername extracts username from context (set by BearerAuthenticator)
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetAuthorities extracts authorities from context (set by BearerAuthenticator)
func GetAuthorities(c *gin.Context) []string {
	authorities, exists := c.Get("authorities")
	if !exists {
		return []string{}
	}
	return authorities.([]string)
}

func extractClaims(tokenStr string, hmacSecret []byte) (*model.Claims, bool) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	claims := &model.Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		log.Printf("Invalid JWT Token: %v", err)
		return nil, false
	}
	// Check if token matches the latest token stored in Redis (ensure it's the current valid token)
	storedToken, err := redis.Get(redis.Rdb, "auth:token:"+claims.Id)
	if err != nil {
		log.Printf("Error getting token from Redis: %v", err)
		return nil, false
	}
	if storedToken == "" || storedToken != tokenStr {
		log.Printf("Token does not match the latest token in Redis or has been invalidated")
		return nil, false
	}
	return claims, true
}
