package service

import (
	"context"

	"github.com/friedrichad/golang_web_api_demo/internal/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
)

type IPositionLevelService interface {
	CacheLevel(ctx context.Context, userID, level int, position string) error
	GetLevel(ctx context.Context, userID int) int
	InvalidateLevel(ctx context.Context, userID int) error
	CanManage(userLevel, targetLevel int) bool
}

type PositionLevelService struct {
	manager *utils.PositionLevelManager
}

func NewPositionLevelService() IPositionLevelService {
	return &PositionLevelService{
		manager: utils.NewPositionLevelManager(redis.Rdb),
	}
}

func (ps *PositionLevelService) CacheLevel(ctx context.Context, userID, level int, position string) error {
	return ps.manager.CacheLevel(ctx, userID, level, position)
}

func (ps *PositionLevelService) GetLevel(ctx context.Context, userID int) int {
	return ps.manager.GetLevel(ctx, userID)
}

func (ps *PositionLevelService) InvalidateLevel(ctx context.Context, userID int) error {
	ps.manager.InvalidateLevel(ctx, userID)
	return nil
}

func (ps *PositionLevelService) CanManage(userLevel, targetLevel int) bool {
	return ps.manager.CanManage(userLevel, targetLevel)
}
