package dtos

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// InventoryAdjustmentDetailFilter - GET request with query parameters
type InventoryAdjustmentDetailFilter struct {
	AdjustmentDetailID int `form:"adjustment_detail_id"`
	AdjustmentID       int `form:"adjustment_id"`
	model.PageSize
	model.DateRequest
}

// InventoryAdjustmentDetailCreate - POST request body
type InventoryAdjustmentDetailCreate struct {
	AdjustmentID       int     `json:"adjustment_id" binding:"required"`
	ComponentID        int     `json:"component_id" binding:"required"`
	BinID              int     `json:"bin_id" binding:"required"`
	WarehouseID        int     `json:"warehouse_id" binding:"required"`
	QuantityBefore     float64 `json:"quantity_before" binding:"required"`
	QuantityAfter      float64 `json:"quantity_after" binding:"required"`
	AdjustmentQuantity float64 `json:"adjustment_quantity" binding:"required"`
}

// InventoryAdjustmentDetailUpdate - PUT request body
type InventoryAdjustmentDetailUpdate struct {
	AdjustmentDetailID int     `json:"adjustment_detail_id" binding:"required"`
	AdjustmentID       int     `json:"adjustment_id"`
	ComponentID        int     `json:"component_id"`
	BinID              int     `json:"bin_id"`
	WarehouseID        int     `json:"warehouse_id"`
	QuantityBefore     float64 `json:"quantity_before"`
	QuantityAfter      float64 `json:"quantity_after"`
	AdjustmentQuantity float64 `json:"adjustment_quantity"`
	UpdatedBy          int     `json:"updated_by"`
}

type InventoryAdjustmentDetailResponse struct {
	AdjustmentDetailID int     `json:"adjustment_detail_id"`
	AdjustmentID       int     `json:"adjustment_id"`
	ComponentID        int     `json:"component_id"`
	BinID              int     `json:"bin_id"`
	WarehouseID        int     `json:"warehouse_id"`
	QuantityBefore     float64 `json:"quantity_before"`
	QuantityAfter      float64 `json:"quantity_after"`
	AdjustmentQuantity float64 `json:"adjustment_quantity"`
}
