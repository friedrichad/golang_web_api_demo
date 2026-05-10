package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
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
		jwtClaims, valid := extractClaims(c.GetHeader("Authorization"), hmacSecret)
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
		c.Next()
	}
}

// Authorizator checks if user has required authorities/permissions
// Can be called with multiple authorities - user needs at least one of them
func Authorizator(authority ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities := c.GetStringSlice("authorities")
		if len(authorities) == 0 || !utils.AnyContains(authorities, authority) {
			c.JSON(http.StatusForbidden, model.ResponseWrapper{
				Code:    "403",
				Message: "Không có quyền truy cập",
			})
			c.Abort()
			return
		}
		c.Next()
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
	if err == nil {
		return claims, true
	}
	log.Printf("Invalid JWT Token: %v", err)
	return nil, false
}
