package dtos

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type DepartmentFilter struct {
	DepartmentID   int       `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	Description    string    `json:"description"`
	model.PageSize
	model.DateRequest
}

type DepartmentCreate struct {
	DepartmentName string `json:"department_name" binding:"required"`
	Description    string `json:"description"`
}

type DepartmentUpdate struct {
	DepartmentID   int    `json:"department_id" binding:"required"`
	DepartmentName string `json:"department_name"`
	Description    string `json:"description"`
}

type DepartmentResponse struct {
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	Description    string `json:"description"`
}