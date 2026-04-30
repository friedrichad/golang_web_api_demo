package dtos

import "github.com/friedrichad/golang_web_api_demo/internal/model"

type RoleDTO struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

// RoleFilter - GET request with query parameters
type RoleFilter struct {
	RoleName    string `form:"role_name"`
	Description string `form:"description"`
	model.PageSize
	model.DateRequest
}

// RoleCreate - POST request body
type RoleCreate struct {
	RoleName    string `json:"role_name" binding:"required"`
	Description string `json:"description"`
}

// RoleUpdate - PUT request body
type RoleUpdate struct {
	RoleID      int    `json:"role_id" binding:"required"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	UpdatedBy   int    `json:"updated_by"`
}

type RoleResponse struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}
