package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

func SaveUserPermissionCache(rdb *redis.Client, userId int, scopes []string, ttl time.Duration) error {
	if len(scopes) == 0 {
		return nil
	}
	key := fmt.Sprintf("user_permission:%d", userId)
	pipe := rdb.Pipeline()
	ttlSec := int(ttl.Seconds())
	for _, scope := range scopes {
		pipe.HSet(Ctx, key, scope, 1)
		pipe.Do(Ctx, "HEXPIRE", key, ttlSec, "FIELDS", 1, scope)
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

// CheckRestrictedMenuPermission checks if a user has permission for a restricted menu
// Returns true only if:
// 1. The scope is in the restricted permissions list
// 2. The user has that scope in their user_permission hash
func CheckRestrictedMenuPermission(rdb *redis.Client, userId int, scope string) bool {
	// Get restricted permissions list
	restrictedListKey := "restricted_permissions:list"
	restrictedListData, err := Get(rdb, restrictedListKey)
	if err != nil || restrictedListData == "" {
		log.Printf("[AUTH] Failed to get restricted permissions list: %v", err)
		return false
	}

	// Check if scope is in restricted list
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

	// Check if user has this permission in their hash
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
