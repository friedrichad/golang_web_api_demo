package middleware

import (
	"errors"
	"net/http"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
)

func CheckPositionLevel(requiredLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		levelInterface, exists := c.Get("position_level")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ResponseWrapper{
				Code:    "401",
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		userLevel := utils.GetUserLevelFromContext(levelInterface)
		if userLevel > requiredLevel {
			c.JSON(http.StatusForbidden, model.ResponseWrapper{
				Code:    "403",
				Message: "Access denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func VerifyCanManageUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestorLevelInterface, exists := c.Get("position_level")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ResponseWrapper{
				Code:    "401",
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		requestorLevel := utils.GetUserLevelFromContext(requestorLevelInterface)
		targetUserIDStr := c.Param("user_id")

		if targetUserIDStr == "" {
			c.Next()
			return
		}

		targetLevel, err := getTargetUserLevel(c, targetUserIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ResponseWrapper{
				Code:    "400",
				Message: "User not found",
			})
			c.Abort()
			return
		}

		if requestorLevel >= targetLevel {
			c.JSON(http.StatusForbidden, model.ResponseWrapper{
				Code:    "403",
				Message: "Access denied",
			})
			c.Abort()
			return
		}

		c.Set("target_level", targetLevel)
		c.Next()
	}
}

func getTargetUserLevel(c *gin.Context, userIDStr string) (int, error) {
	ctx := c.Request.Context()
	manager := utils.NewPositionLevelManager(redis.Rdb)

	userID, err := utils.StringToInt(userIDStr)
	if err != nil {
		return 0, err
	}

	level := manager.GetLevel(ctx, userID)
	if level != 999 {
		return level, nil
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetById(userID)
	if err != nil || user == nil {
		return 0, errors.New("user not found")
	}

	manager.CacheLevel(ctx, userID, user.PositionLevel, user.PositionName)

	return user.PositionLevel, nil
}

