package dtos

import "github.com/friedrichad/golang_web_api_demo/internal/model"

type RoleDTO struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

type RoleRequest struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	model.PageSize
	model.DateRequest
}
type RoleResponse struct {
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}
