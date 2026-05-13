package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type RequestPermissionFilter struct {
	RequestPermissionID int       `json:"request_permission_id"`
	RequestID           int       `json:"request_id"`
	MenuID              int       `json:"menu_id"`
	PermissionID        int       `json:"permission_id"`
	model.PageSize
	model.DateRequest	
}

type RequestPermissionCreate struct {
	RequestID           int       `json:"request_id"`
	MenuID              int       `json:"menu_id"`
	PermissionID        int       `json:"permission_id"`
	Reason              string    `json:"reason"`
	CreatedAt           time.Time `json:"created_at"`
}

type RequestPermissionUpdate struct {
	RequestPermissionID int       `json:"request_permission_id"`
	RequestID           int       `json:"request_id"`
	MenuID              int       `json:"menu_id"`
	PermissionID        int       `json:"permission_id"`
	Reason              string    `json:"reason"`
}
type RequestPermissionResponse struct {
	RequestPermissionID int       `json:"request_permission_id"`
	RequestID           int       `json:"request_id"`
	MenuID              int       `json:"menu_id"`
	PermissionID        int       `json:"permission_id"`
	Reason              string    `json:"reason"`
	CreatedAt           time.Time `json:"created_at"`
}