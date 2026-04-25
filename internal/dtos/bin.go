package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type BinRequest struct {
	BinID               int       `json:"bin_id"`
	LocationInWarehouse string    `json:"location_in_warehouse"`
	StatusInt           int       `json:"status_int"`
	WarehouseID         int       `json:"warehouse_id"`
	CreatedBy           int       `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	model.PageSize      `json:"page_size"`
	model.DateRequest   `json:"date_request"`
}
type BinResponse struct {
	BinID               int       `json:"bin_id"`
	LocationInWarehouse string    `json:"location_in_warehouse"`
	StatusInt           int       `json:"status_int"`
	WarehouseID         int       `json:"warehouse_id"`
	CreatedBy           int       `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedBy           int       `json:"updated_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	model.PageSize      `json:"page_size"`
	model.DateRequest   `json:"date_request"`
}
