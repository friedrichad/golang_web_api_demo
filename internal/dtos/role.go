package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type RoleDTO struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

// RoleFilter - GET request with query parameters
type RoleFilter struct {
	RoleName *string `form:"role_name"`
	model.PageSize
	model.DateRequest
}

// RoleCreate - POST request body
type RoleCreate struct {
	RoleName    string `json:"role_name" binding:"required"`
	Description string `json:"description"`
}

// Verify validates the RoleCreate struct.
func (r *RoleCreate) Verify() error {
	if r.RoleName == "" {
		return fmt.Errorf("RoleName is required")
	}
	return nil
}

// RoleUpdate - PUT request body
type RoleUpdate struct {
	RoleID      int    `json:"role_id" binding:"required"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	UpdatedBy   int    `json:"updated_by"`
}

// Verify validates the RoleUpdate struct.
func (r *RoleUpdate) Verify() error {
	if r.RoleID == 0 {
		return fmt.Errorf("RoleID is required")
	}
	return nil
}

type RoleResponse struct {
	RoleID      int       `json:"role_id"`
	RoleName    string    `json:"role_name"`
	Description string    `json:"description"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleMenuAssign struct {
	RoleID  int   `json:"role_id" binding:"required"`
	MenuIDs []int `json:"menu_ids"`
}

type PermissionDTO struct {
	PermissionID   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}

type PermissionCreate struct {
	PermissionName string `json:"permission_name" binding:"required"`
}

type PermissionUpdate struct {
	PermissionID   int    `json:"permission_id" binding:"required"`
	PermissionName string `json:"permission_name" binding:"required"`
}
