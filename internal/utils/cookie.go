package utils

import (
	"time"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetOrCreateSessionID(c *gin.Context, ttl time.Duration) string {
	sessionID, err := c.Cookie("session_id")

	if err != nil || sessionID == "" {
		sessionID = uuid.New().String()

		c.SetCookie(
			"session_id",
			sessionID,
			int(ttl.Seconds()), // TTL in seconds
			"/",
			"",
			false,
			true,
		)
	}

	return sessionID
}

func SaveBrowserSession(sessionID string, userID string, ttl time.Duration) error {
	key := "auth:browser:" + sessionID
	return redis.Save(redis.Rdb, key, userID, ttl)
}

func BrowserHasSession(sessionID string) (string, error) {
	key := "auth:browser:" + sessionID
	return redis.Get(redis.Rdb, key)
}

func ClearSessionCookie(c *gin.Context) {
	c.SetCookie(
		"session_id",
		"",
		-1,   // MaxAge < 0 -> delete cookie
		"/",
		"",
		false,
		true,
	)
}