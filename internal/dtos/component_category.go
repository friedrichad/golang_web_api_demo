package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type ComponentCategoryFilter struct {
	CategoryID   int    `form:"category_id"`
	CategoryName string `form:"category_name"`
	model.PageSize
}

type ComponentCategoryCreate struct {
	CategoryName string `json:"category_name" binding:"required"`
}

type ComponentCategoryUpdate struct {
	CategoryID   int    `json:"category_id" binding:"required"`
	CategoryName string `json:"category_name"`
	UpdatedBy    int    `json:"updated_by"`
}

type ComponentCategoryResponse struct {
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    int       `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}
