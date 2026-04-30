package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type UserResponse struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	StatusInt   int       `json:"status_int"`
	UserRole    []string  `json:"user_role"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserFilter - GET request with query parameters
type UserFilter struct {
	Username    string `form:"username"`
	DisplayName string `form:"display_name"`
	Email       string `form:"email"`
	StatusInt   int    `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// UserCreate - POST request body
type UserCreate struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	RoleIDs     []int  `json:"role_ids"`
}

// UserUpdate - PUT request body
type UserUpdate struct {
	UserID      int32 `json:"user_id" binding:"required"`
	Username    string `json:"username" binding:"min=3,max=50"`
	DisplayName string `json:"display_name" binding:"min=1,max=100"`
	Email       string `json:"email" binding:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" binding:"min=6"`
	StatusInt   int32   `json:"status_int"`
	UpdatedBy   int32   `json:"updated_by"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
