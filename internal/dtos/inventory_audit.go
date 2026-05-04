package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// InventoryAuditFilter - GET request with query parameters
type InventoryAuditFilter struct {
	AuditID     *int `form:"audit_id"`
	WarehouseID *int `form:"warehouse_id"`
	StatusInt   *int `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// InventoryAuditCreate - POST request body
type InventoryAuditCreate struct {
	WarehouseID int    `json:"warehouse_id" binding:"required"`
	Note        string `json:"note"`
}

// Verify validates the InventoryAuditCreate struct.
func (i *InventoryAuditCreate) Verify() error {
	if i.WarehouseID == 0 {
		return fmt.Errorf("WarehouseID is required")
	}
	return nil
}

// InventoryAuditUpdate - PUT request body
type InventoryAuditUpdate struct {
	AuditID   int    `json:"audit_id" binding:"required"`
	StatusInt int    `json:"status_int"`
	Note      string `json:"note"`
	UpdatedBy int    `json:"updated_by"`
}

// Verify validates the InventoryAuditUpdate struct.
func (i *InventoryAuditUpdate) Verify() error {
	if i.AuditID == 0 {
		return fmt.Errorf("AuditID is required")
	}
	return nil
}

type InventoryAuditResponse struct {
	AuditID     int       `json:"audit_id"`
	WarehouseID int       `json:"warehouse_id"`
	StatusInt   int       `json:"status_int"`
	Note        string    `json:"note"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}
