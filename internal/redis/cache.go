package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/shared"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		log.Printf("Không kết nối được Redis: %v", err)
		return
	}
	log.Printf("Kết nối thành công với Redis")
}

func Save(rdb *redis.Client, key string, value string, expiration time.Duration) error {
	return rdb.Set(Ctx, key, value, expiration).Err()
}

func Get(rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	return val, err
}

func Delete(rdb *redis.Client, key string) error {
	return rdb.Del(Ctx, key).Err()
}

func Exists(rdb *redis.Client, key string) (bool, error) {
	n, err := rdb.Exists(Ctx, key).Result()
	return n > 0, err
}

func UpdateTTL(rdb *redis.Client, key string, expiration time.Duration) error {
	return rdb.Expire(Ctx, key, expiration).Err()
}

func GetTTL(rdb *redis.Client, key string) (time.Duration, error) {
	return rdb.TTL(Ctx, key).Result()
}

func AddToBlacklist(token string, ttl time.Duration) error {
	key := "blacklist:" + token
	return Rdb.Set(Ctx, key, "1", ttl).Err()
}

func IsBlacklisted(token string) (bool, error) {
	key := "blacklist:" + token
	n, err := Rdb.Exists(Ctx, key).Result()
	return n > 0, err
}

func SaveUserPermissionCache(rdb *redis.Client, userId int, scopes []shared.UserPermissionScope) error {
	if len(scopes) == 0 {
		return nil
	}
	key := fmt.Sprintf("user_permission:%d", userId)
	pipe := rdb.Pipeline()
	for _, scope := range scopes {
		scope_ttl := time.Until(time.Unix(scope.ExpiredDate, 0))
		ttlSec := int(scope_ttl.Seconds())
		pipe.HSet(Ctx, key, scope.Scope, 1)
		pipe.Do(Ctx, "HEXPIRE", key, ttlSec, "FIELDS", 1, scope.Scope)
	}
	_, err := pipe.Exec(Ctx)
	return err
}

func CheckPermissionRedis(rdb *redis.Client, userId int, authorities []string) bool {
	key := fmt.Sprintf("user_permission:%d", userId)
	pipe := rdb.Pipeline()
	cmds := make([]*redis.BoolCmd, len(authorities))
	for i, scope := range authorities {
		cmds[i] = pipe.HExists(Ctx, key, scope)
	}
	_, err := pipe.Exec(Ctx)
	if err != nil {
		log.Printf("Redis pipeline error: %v", err)
		return false
	}
	for _, cmd := range cmds {
		if ok, _ := cmd.Result(); ok {
			return true
		}
	}
	return false
}

// 1. The scope is in the restricted permissions list
// 2. The user has that scope in their user_permission hash
func CheckRestrictedMenuPermission(rdb *redis.Client, userId int, scope string) bool {
	restrictedListKey := "restricted_permissions:list"
	restrictedListData, err := Get(rdb, restrictedListKey)
	if err != nil || restrictedListData == "" {
		log.Printf("[AUTH] Failed to get restricted permissions list: %v", err)
		return false
	}

	var restrictedScopes []string
	if err := json.Unmarshal([]byte(restrictedListData), &restrictedScopes); err != nil {
		log.Printf("[AUTH] Failed to parse restricted scopes: %v", err)
		return false
	}
	isRestricted := false
	for _, r := range restrictedScopes {
		if r == scope {
			isRestricted = true
			break
		}
	}
	if !isRestricted {
		log.Printf("[AUTH] Scope '%s' is not in restricted list", scope)
		return false
	}
	userPermKey := fmt.Sprintf("user_permission:%d", userId)
	hasPermission, err := rdb.HExists(Ctx, userPermKey, scope).Result()
	if err != nil {
		log.Printf("[AUTH] Error checking user permission hash: %v", err)
		return false
	}
	if !hasPermission {
		log.Printf("[AUTH] User %d does not have permission '%s'", userId, scope)
		return false
	}

	log.Printf("[AUTH] ✓ User %d has restricted menu permission '%s'", userId, scope)
	return true
}

func DeleteUserPermissionField(rdb *redis.Client, userId int) error {
	key := fmt.Sprintf("user_permission:%d", userId)
	return rdb.Del(Ctx, key).Err()
}

func SaveUserInfoCache(rdb *redis.Client, userInfo shared.UserInfo, expiration time.Duration) error {
	key := fmt.Sprintf("user_info:%d", userInfo.UserId)
	data, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}
	return Save(rdb, key, string(data), expiration)
}

func DeleteUserInfoCache(rdb *redis.Client, userId int) error {
	key := fmt.Sprintf("user_info:%d", userId)
	return rdb.Del(Ctx, key).Err()
}

func CanApproveRequest(rdb *redis.Client, approverId int, requesterPositionLevel int) (bool, error) {
	cacheKey := fmt.Sprintf("user_info:%d", approverId)
	approverCacheData, err := Get(rdb, cacheKey)
	if err != nil {
		log.Printf("[APPROVAL] Failed to get approver cache: %v", err)
		return false, err
	}
	var approverInfo shared.UserInfo
	if approverCacheData == "" {
		log.Printf("[APPROVAL] Approver ID=%d not in cache, must load from database", approverId)
		return false, fmt.Errorf("approver not found in cache")
	}
	if err := json.Unmarshal([]byte(approverCacheData), &approverInfo); err != nil {
		log.Printf("[APPROVAL] Failed to parse approver cache data: %v", err)
		return false, err
	}

	approverLevel := approverInfo.PositionInfo.PositionLevel
	canApprove := approverLevel >= requesterPositionLevel

	if canApprove {
		log.Printf("[APPROVAL] ✓ Approver ID=%d (Level=%d) CAN approve requester (Level=%d)",
			approverId, approverLevel, requesterPositionLevel)
	} else {
		log.Printf("[APPROVAL] ✗ Approver ID=%d (Level=%d) CANNOT approve requester (Level=%d) - insufficient position level",
			approverId, approverLevel, requesterPositionLevel)
	}

	return canApprove, nil
}
func CompareUserInfoCache(rdb *redis.Client, currentUserId int, targetUser shared.UserInfo) (bool, error) {
	log.Printf("[DEPRECATION] CompareUserInfoCache is deprecated, use CanApproveRequest instead")
	return CanApproveRequest(rdb, currentUserId, targetUser.PositionInfo.PositionLevel)
}
