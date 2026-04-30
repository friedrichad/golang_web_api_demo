package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)
type WarehouseRequest struct{
	WarehouseID      int       `json:"warehouse_id"`
	WarehouseName    string    `json:"warehouse_name"`
	ImageURL         string    `json:"image_url"`
	CreatedBy        int       `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedBy        int       `json:"updated_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	model.PageSize
	model.DateRequest
}
type WarehouseResponse struct{
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
type WarehouseUpdate struct{
	WarehouseID      int       `json:"warehouse_id"`
	WarehouseName    string    `json:"warehouse_name"`
	Description      string    `json:"description"`
	PhysicalLocation string    `json:"physical_location"`
	ImageURL         string    `json:"image_url"`
}