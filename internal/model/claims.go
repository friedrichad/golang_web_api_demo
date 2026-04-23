package model

import (
	"errors"
	"time"
)

type Claims struct {
	Id          string   `json:"id"`
	Username    string   `json:"user_name"`
	Authorities []string `json:"authorities"`
	ClientId    string   `json:"client_id"`
	Exp         int64    `json:"exp"`
	RefreshExp  int64    `json:"refresh_exp"`
}

func (c Claims) Valid() error {
	if time.Unix(0, c.Exp*int64(time.Second)).Before(time.Now()) {
		return errors.New("Token is expired")
	}
	return nil
}

func (c Claims) RefreshTokenExpired() bool {
	if time.Unix(0, c.RefreshExp*int64(time.Second)).Before(time.Now()) {
		return false
	}
	return true
}
