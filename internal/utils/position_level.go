package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PositionLevelManager struct {
	rdb *redis.Client
}

func NewPositionLevelManager(rdb *redis.Client) *PositionLevelManager {
	return &PositionLevelManager{rdb: rdb}
}

func cacheKeyUserLevel(userID int) string {
	return fmt.Sprintf("position:level:user:%d", userID)
}

type UserLevelInfo struct {
	UserID   int    `json:"user_id"`
	Level    int    `json:"level"`
	Position string `json:"position"`
}

func (plm *PositionLevelManager) CacheLevel(ctx context.Context, userID, level int, position string) error {
	info := UserLevelInfo{UserID: userID, Level: level, Position: position}
	data, _ := json.Marshal(info)
	return plm.rdb.Set(ctx, cacheKeyUserLevel(userID), string(data), 1*time.Hour).Err()
}

func (plm *PositionLevelManager) GetLevel(ctx context.Context, userID int) int {
	val, err := plm.rdb.Get(ctx, cacheKeyUserLevel(userID)).Result()
	if err != nil {
		return 999
	}
	var info UserLevelInfo
	json.Unmarshal([]byte(val), &info)
	return info.Level
}

func (plm *PositionLevelManager) InvalidateLevel(ctx context.Context, userID int) {
	plm.rdb.Del(ctx, cacheKeyUserLevel(userID))
}

func (plm *PositionLevelManager) CanManage(userLevel, targetLevel int) bool {
	return userLevel < targetLevel
}

func GetUserLevelFromContext(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	default:
		return 999
	}
}
