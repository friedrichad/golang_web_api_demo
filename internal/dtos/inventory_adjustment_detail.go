package dtos

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)
type InventoryAdjustmentDetailRequest struct{
	AdjustmentDetailID int     `json:"adjustment_detail_id"`
	AdjustmentID       int     `json:"adjustment_id"`
	model.PageSize
	model.DateRequest
}
type InventoryAdjustmentDetailResponse struct{
	AdjustmentDetailID int     `json:"adjustment_detail_id"`
	AdjustmentID       int     `json:"adjustment_id"`
	ComponentID        int     `json:"component_id"`
	BinID              int     `json:"bin_id"`
	WarehouseID        int     `json:"warehouse_id"`
	QuantityBefore     float64 `json:"quantity_before"`
	QuantityAfter      float64 `json:"quantity_after"`
	AdjustmentQuantity float64 `json:"adjustment_quantity"`
}
type InventoryAdjustmentDetailUpdate struct{
	AdjustmentDetailID int     `json:"adjustment_detail_id"`
	AdjustmentID       int     `json:"adjustment_id"`
	ComponentID        int     `json:"component_id"`
	BinID              int     `json:"bin_id"`
	WarehouseID        int     `json:"warehouse_id"`
	QuantityBefore     float64 `json:"quantity_before"`
	QuantityAfter      float64 `json:"quantity_after"`
	AdjustmentQuantity float64 `json:"adjustment_quantity"`
}