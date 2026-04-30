package dtos

import (
	"time"
)

type SystemLogDTO struct {
	LogID      int       `json:"log_id"`
	UserID     int       `json:"user_id"`
	HTTPMethod string    `json:"http_method"`
	Route      string    `json:"route"`
	StatusInt  int       `json:"status_int"`
	IPAddress  string    `json:"ip_address"`
	ExecutedAt time.Time `json:"executed_at"`
}