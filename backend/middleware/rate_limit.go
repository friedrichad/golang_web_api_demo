package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/friedrichad/golang_web_api_demo/backend/redis"
	"github.com/gin-gonic/gin"
)

func RateLimit(prefix string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := prefix + ":" + ip
		count, err := redis.Rdb.Incr(context.Background(), key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			redis.Rdb.Expire(context.Background(), key, window)
		}
		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests",
			},
			)
			return
		}
		c.Next()
	}
}
