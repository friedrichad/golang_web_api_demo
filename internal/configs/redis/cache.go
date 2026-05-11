package redis

import (
	"context"
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


// Create hoặc Update (Set)
func Save(rdb *redis.Client, key string, value string, expiration time.Duration) error {
    return rdb.Set(Ctx, key, value, expiration).Err()
}

// Read (Get)
func Get(rdb *redis.Client, key string) (string, error) {
    return rdb.Get(Ctx, key).Result()
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