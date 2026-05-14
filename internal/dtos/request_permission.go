package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type RequestPermissionFilter struct {
	RequestPermissionID int       `form:"request_permission_id"`
	RequestID           int       `form:"request_id"`
	MenuID              int       `form:"menu_id"`
	PermissionID        int       `form:"permission_id"`
	model.PageSize
	model.DateRequest	
}

type RequestPermissionCreate struct {
	RequestID           int       `json:"request_id", binding:"required"`
	MenuID              int       `json:"menu_id", binding:"required"`
	PermissionID        int       `json:"permission_id", binding:"required"`
	Reason              string    `json:"reason", binding:"required"`
}

type RequestPermissionUpdate struct {
	RequestPermissionID int       `json:"request_permission_id" binding:"required"`
	RequestID           int       `json:"request_id" binding:"required"`
	MenuID              int       `json:"menu_id" binding:"required"`
	PermissionID        int       `json:"permission_id" binding:"required"`
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