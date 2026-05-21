package middleware

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// TTL constants for Redis cache
const (
	UserPermissionCacheTTL       = 24 * time.Hour // User permissions cache TTL
	RestrictedPermissionCacheTTL = 1 * time.Hour  // Restricted permissions cache TTL
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
		rawToken := strings.Replace(authHeader, "Bearer ", "", 1)
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
// Check flow:
// 1. Check if operator -> bypass all checks
// 2. Check restricted permissions first (user must have sufficient position_level)
// 3. Check JWT authorities
// 4. Check Redis cache
// 5. Load from DB and verify (lazy load + cache)
func Authorizator(authority ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isOperator(c) {
			c.Next()
			return
		}
		userId, err := strconv.Atoi(GetUserID(c))
		if err != nil {
			forbidden(c)
			return
		}
		if hasRestrictedAuthority(authority) {
			if !checkRestrictedPermissions(c, userId, authority) {
				forbidden(c)
				return
			}
			c.Next()
			return
		}
		allowed :=
			checkJWTAuthorities(c, authority) ||
				checkRedisPermission(userId, authority) ||
				loadAndCheckUserPermissions(userId, authority)

		if allowed {
			c.Next()
			return
		}

		forbidden(c)
	}
}
func forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, model.ResponseWrapper{
		Code:    "403",
		Message: "Không có quyền truy cập, vui lòng liên hệ quản trị viên",
	})
	c.Abort()
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

func isOperator(c *gin.Context) bool {
	isOp := c.GetInt("is_op") == 1
	if isOp {
		log.Printf("[AUTH] User is operator - bypassing permission checks")
	}
	return isOp
}

func checkJWTAuthorities(c *gin.Context, authority []string) bool {
	jwtAuthorities := c.GetStringSlice("authorities")
	if len(jwtAuthorities) == 0 {
		log.Printf("[AUTH] JWT: No authorities in token claims")
		return false
	}

	if !utils.AnyContains(jwtAuthorities, authority) {
		log.Printf("[AUTH] JWT: Check failed - JWT: %v, Required: %v", jwtAuthorities, authority)
		return false
	}

	log.Printf("[AUTH] JWT: ✓ User has required authority")
	return true
}

// checkRedisPermission checks if user has required authorities in Redis cache
func checkRedisPermission(userId int, authority []string) bool {
	if !redis.CheckPermissionRedis(redis.Rdb, userId, authority) {
		log.Printf("[AUTH] Redis: ✗ User %d - authorities %v not in cache", userId, authority)
		return false
	}

	log.Printf("[AUTH] Redis: ✓ User %d has required authority from cache", userId)
	return true
}

// hasRestrictedAuthority checks if any requested authority is a restricted permission
func hasRestrictedAuthority(authorities []string) bool {
	restrictedPerms, err := getRestrictedPermissionsList()
	if err != nil {
		log.Printf("[AUTH] ⚠ Cannot load restricted permissions list: %v", err)
		return false
	}

	for _, auth := range authorities {
		if isPermissionRestricted(auth, restrictedPerms) {
			log.Printf("[AUTH] Detected restricted authority: %s", auth)
			return true
		}
	}
	return false
}

func checkRestrictedPermissions(c *gin.Context, userId int, authorities []string) bool {
	log.Printf("[AUTH] RESTRICTED: Checking for user_id %d", userId)
	positionLevel := c.GetInt("position_level")
	if positionLevel <= 0 {
		log.Printf("[AUTH] RESTRICTED: ✗ User position_level %d insufficient", positionLevel)
		return false
	}
	log.Printf("[AUTH] RESTRICTED: ✓ User position_level %d valid", positionLevel)

	// Tối ưu: Kiểm tra quyền trong Redis Cache trước khi gọi Database
	if checkRedisPermission(userId, authorities) {
		log.Printf("[AUTH] RESTRICTED: ✓ User %d authorized from cache", userId)
		return true
	}
	userRepo := repository.NewUserRepository()
	perms, err := userRepo.GetUserPermissionScopes(userId)
	if err != nil {
		log.Printf("[AUTH] RESTRICTED: ✗ Failed to load permissions: %v", err)
		return false
	}

	if len(perms) == 0 {
		log.Printf("[AUTH] RESTRICTED: ✗ User %d has no permissions", userId)
		return false
	}

	scopes := make([]string, len(perms))
	for i, p := range perms {
		scopes[i] = p.Scope
	}

	// Verify authority
	if !checkUserHasAuthority(scopes, authorities) {
		log.Printf("[AUTH] RESTRICTED: ✗ User %d - Required: %v, Has: %v", userId, authorities, scopes)
		return false
	}

	// Cache permissions
	err = redis.SaveUserPermissionCache(redis.Rdb, userId, perms)
	if err != nil {
		log.Printf("[AUTH-CACHE] ⚠ Failed to cache: %v", err)
	} else {
		log.Printf("[AUTH-CACHE] ✓ Cached %d permissions", len(perms))
	}

	log.Printf("[AUTH] RESTRICTED: ✓ User %d authorized", userId)
	return true
}

func loadAndCheckUserPermissions(userId int, authorities []string) bool {
	log.Printf("[AUTH] DB: Loading permissions for user_id %d", userId)

	userRepo := repository.NewUserRepository()
	perms, err := userRepo.GetUserPermissionScopes(userId)
	if err != nil {
		log.Printf("[AUTH] DB: ✗ Failed to load - %v", err)
		return false
	}

	if len(perms) == 0 {
		log.Printf("[AUTH] DB: ✗ User %d has no permissions", userId)
		return false
	}

	// Convert to scopes
	scopes := make([]string, len(perms))
	for i, p := range perms {
		scopes[i] = p.Scope
	}

	// Check authority
	if !checkUserHasAuthority(scopes, authorities) {
		log.Printf("[AUTH] DB: ✗ User %d - Required: %v, Has: %v", userId, authorities, scopes)
		return false
	}

	// Cache for future use
	err = redis.SaveUserPermissionCache(redis.Rdb, userId, perms)
	if err != nil {
		log.Printf("[AUTH-CACHE] ⚠ Failed to cache for user %d: %v", userId, err)
	} else {
		log.Printf("[AUTH-CACHE] ✓ Cached %d permissions for user %d (TTL: %v)", len(perms), userId, UserPermissionCacheTTL)
	}

	log.Printf("[AUTH] DB: ✓ User %d has required authority", userId)
	return true
}

func checkUserHasAuthority(userScopes []string, requiredAuthorities []string) bool {
	for _, required := range requiredAuthorities {
		for _, userScope := range userScopes {
			if userScope == required {
				return true
			}
		}
	}
	return false
}

func InitUserPermissionCache(userID int) bool {
	userRepo := repository.NewUserRepository()
	perms, err := userRepo.GetUserPermissionScopes(userID)
	if err != nil {
		log.Printf("Failed to get user permissions: %v", err)
		return false
	}
	if len(perms) == 0 {
		log.Printf("User %d has no permissions to cache", userID)
		return false
	}

	scopes := make([]string, len(perms))
	for i, p := range perms {
		scopes[i] = p.Scope
	}

	// Cache entire list first
	err = redis.SaveUserPermissionCache(redis.Rdb, userID, perms)
	if err != nil {
		log.Printf("[AUTH-CACHE] ✗ Failed to cache user permission list for user_id %d (TTL: %v): %v", userID, UserPermissionCacheTTL, err)
		return false
	}
	log.Printf("[AUTH-CACHE] ✓ Cached %d permissions for user %d (TTL: %v)", len(perms), userID, UserPermissionCacheTTL)

	return true
}

func getRestrictedPermissionsList() ([]string, error) {
	cacheKey := "restricted_permissions:list"

	// Try to get from Redis cache first
	cachedData, err := redis.Get(redis.Rdb, cacheKey)
	if err == nil && cachedData != "" {
		var scopes []string
		if err := json.Unmarshal([]byte(cachedData), &scopes); err == nil {
			log.Printf("[AUTH-CACHE] ✓ Loaded %d restricted permissions from cache (TTL: %v)", len(scopes), RestrictedPermissionCacheTTL)
			return scopes, nil
		}
	}

	return loadRestrictedPermissionsFromDB(cacheKey)
}

func loadRestrictedPermissionsFromDB(cacheKey string) ([]string, error) {
	menuPermRepo := repository.NewMenuPermissionRepository()

	scopes, err := menuPermRepo.GetAllMenuPermissionsByRestricted(1)
	if err != nil {
		log.Printf("[AUTH-CACHE] ✗ Failed to load restricted permissions from DB: %v", err)
		return []string{}, err
	}

	if len(scopes) == 0 {
		log.Printf("[AUTH-CACHE] No restricted permissions found in database")
		return []string{}, nil
	}
	cacheData, _ := json.Marshal(scopes)
	err = redis.Save(redis.Rdb, cacheKey, string(cacheData), RestrictedPermissionCacheTTL)
	if err != nil {
		log.Printf("[AUTH-CACHE] ⚠ Failed to cache permission list (TTL: %v): %v", RestrictedPermissionCacheTTL, err)
	} else {
		log.Printf("[AUTH-CACHE] ✓ Cached %d restricted permissions with TTL: %v", len(scopes), RestrictedPermissionCacheTTL)
	}

	return scopes, nil
}

func isPermissionRestricted(scope string, restrictedScopes []string) bool {
	for _, restricted := range restrictedScopes {
		if restricted == scope {
			return true
		}
	}
	return false
}
