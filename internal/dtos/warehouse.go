package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// WarehouseFilter - GET request with query parameters
type WarehouseFilter struct {
	WarehouseID   int    `form:"warehouse_id"`
	WarehouseName string `form:"warehouse_name"`
	model.PageSize
	model.DateRequest
}

// WarehouseCreate - POST request body
type WarehouseCreate struct {
	WarehouseName    string `json:"warehouse_name" binding:"required"`
	Description      string `json:"description"`
	PhysicalLocation string `json:"physical_location"`
	ImageURL         string `json:"image_url"`
}

// WarehouseUpdate - PUT request body
type WarehouseUpdate struct {
	WarehouseID      int    `json:"warehouse_id" binding:"required"`
	WarehouseName    string `json:"warehouse_name"`
	Description      string `json:"description"`
	PhysicalLocation string `json:"physical_location"`
	ImageURL         string `json:"image_url"`
	UpdatedBy        int    `json:"updated_by"`
}

// WarehouseResponse - Response body
type WarehouseResponse struct {
	WarehouseID      int       `json:"warehouse_id"`
	WarehouseName    string    `json:"warehouse_name"`
	Description      string    `json:"description"`
	PhysicalLocation string    `json:"physical_location"`
	ImageURL         string    `json:"image_url"`
	CreatedBy        int       `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedBy        int       `json:"updated_by"`
	UpdatedAt        time.Time `json:"updated_at"`
}
