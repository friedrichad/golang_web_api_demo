package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// BinFilter - GET request with query parameters
type BinFilter struct {
	BinID               int    `form:"bin_id"`
	LocationInWarehouse string `form:"location_in_warehouse"`
	StatusInt           int    `form:"status_int"`
	WarehouseID         int    `form:"warehouse_id"`
	model.PageSize
	model.DateRequest
}

// BinCreate - POST request body
type BinCreate struct {
	LocationInWarehouse string `json:"location_in_warehouse" binding:"required"`
	WarehouseID         int    `json:"warehouse_id" binding:"required"`
}

// BinUpdate - PUT request body
type BinUpdate struct {
	BinID               int    `json:"bin_id" binding:"required"`
	LocationInWarehouse string `json:"location_in_warehouse"`
	StatusInt           int    `json:"status_int"`
	WarehouseID         int    `json:"warehouse_id"`
	UpdatedBy           int    `json:"updated_by"`
}

// BinResponse - Response body
type BinResponse struct {
	BinID               int       `json:"bin_id"`
	LocationInWarehouse string    `json:"location_in_warehouse"`
	StatusInt           int       `json:"status_int"`
	WarehouseName       string    `json:"warehouse_name"`
	CreatedBy           int       `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedBy           int       `json:"updated_by"`
	UpdatedAt           time.Time `json:"updated_at"`
}
