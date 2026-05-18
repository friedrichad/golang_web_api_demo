package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

// Create hoặc Update (Set)
func Save(rdb *redis.Client, key string, value string, expiration time.Duration) error {
	return rdb.Set(Ctx, key, value, expiration).Err()
}

// Read (Get)
// Returns empty string if key not found (no error)
func Get(rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(Ctx, key).Result()

	// Handle key not found as normal case (not an error)
	if err == redis.Nil {
		return "", nil
	}

	return val, err
}

// Delete (Del)
func Delete(rdb *redis.Client, key string) error {
	return rdb.Del(Ctx, key).Err()
}

// Exists kiểm tra key có tồn tại không
func Exists(rdb *redis.Client, key string) (bool, error) {
	n, err := rdb.Exists(Ctx, key).Result()
	return n > 0, err
}

// UpdateTTL cập nhật lại thời gian sống của key
func UpdateTTL(rdb *redis.Client, key string, expiration time.Duration) error {
	return rdb.Expire(Ctx, key, expiration).Err()
}

// GetTTL lấy thời gian sống còn lại của key
func GetTTL(rdb *redis.Client, key string) (time.Duration, error) {
	return rdb.TTL(Ctx, key).Result()
}

// AddToBlacklist thêm token vào blacklist với TTL
func AddToBlacklist(token string, ttl time.Duration) error {
	key := "blacklist:" + token
	return Rdb.Set(Ctx, key, "1", ttl).Err()
}

// IsBlacklisted kiểm tra nếu token có trong blacklist
func IsBlacklisted(token string) (bool, error) {
	key := "blacklist:" + token
	n, err := Rdb.Exists(Ctx, key).Result()
	return n > 0, err
}

func SaveUserPermissionCache(rdb *redis.Client, userId int, data map[string]interface{}, ttl time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	key := fmt.Sprintf("user_permission:%d", userId)

	// Convert to strings for Redis and save all at once
	stringData := make(map[string]interface{}, len(data))
	for scope, exp := range data {
		log.Printf("Caching permission for user %d: %s with expiration %v", userId, scope, exp)
		stringData[scope] = fmt.Sprint(exp)
	}

	if err := rdb.HSet(Ctx, key, stringData).Err(); err != nil {
		return err
	}

	return rdb.Expire(Ctx, key, ttl).Err()
}

func CheckPermissionRedis(rdb *redis.Client, userId int, authorities []string) bool {
	key := "user_permission:" + strconv.Itoa(userId)

	// Get only the specific permissions we need to check
	results, err := rdb.HMGet(Ctx, key, authorities...).Result()
	if err != nil || len(results) == 0 {
		return false
	}

	now := time.Now().Unix()
	for _, val := range results {
		if val != nil {
			if exp, err := strconv.ParseInt(val.(string), 10, 64); err == nil {
				if exp == 0 || exp > now {
					return true
				}
			}
		}
	}
	return false
}
