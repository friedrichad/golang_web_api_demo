package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// InventoryLedgerFilter - GET request with query parameters
type InventoryLedgerFilter struct {
	LedgerID        int `form:"ledger_id"`
	ComponentID     int `form:"component_id"`
	WarehouseID     int `form:"warehouse_id"`
	BinID           int `form:"bin_id"`
	ReferenceTypeID int `form:"reference_type_id"`
	model.PageSize
	model.DateRequest
}

type InventoryLedgerResponse struct {
	LedgerID        int       `json:"ledger_id"`
	ComponentID     int       `json:"component_id"`
	WarehouseID     int       `json:"warehouse_id"`
	BinID           int       `json:"bin_id"`
	ReferenceType   int       `json:"reference_type"`
	ReferenceTypeID int       `json:"reference_type_id"`
	Description     string    `json:"description"`
	QuantityChange  float64   `json:"quantity_change"`
	QuantityAfter   float64   `json:"quantity_after"`
	Note            string    `json:"note"`
	CreatedBy       int       `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}
