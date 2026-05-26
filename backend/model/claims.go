package model

import (
	"errors"
	"time"
)

type Claims struct {
	Id           string   `json:"id"`
	Username     string   `json:"user_name"`
	PositionID   int      `json:"position_id"`
	PositionName string   `json:"position_name"`
	IsOP         int      `json:"is_op"`
	Level        int      `json:"position_level"`
	Authorities  []string `json:"authorities"`
	ClientId     string   `json:"client_id"`
	Exp          int64    `json:"exp"`
	RefreshExp   int64    `json:"refresh_exp"`
}

func (c Claims) Valid() error {
	if time.Unix(0, c.Exp*int64(time.Second)).Before(time.Now()) {
		return errors.New("Token is expired")
	}
	return nil
}

// IsRefreshTokenValid checks if the refresh token has NOT expired
func (c Claims) IsRefreshTokenValid() bool {
	return time.Unix(0, c.RefreshExp*int64(time.Second)).After(time.Now())
}

// RefreshTokenExpired checks if refresh token has expired (deprecated - use IsRefreshTokenValid instead)
func (c Claims) RefreshTokenExpired() bool {
	if time.Unix(0, c.RefreshExp*int64(time.Second)).Before(time.Now()) {
		return false
	}
	return true
}
