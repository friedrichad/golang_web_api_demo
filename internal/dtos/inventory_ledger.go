package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)
type InventoryLedgerRequest struct {
	LedgerID        int       `json:"ledger_id"`
	WarehouseID     int       `json:"warehouse_id"`
	ReferenceTypeID int       `json:"reference_type_id"`
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