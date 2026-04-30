package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type InventoryAuditRequest struct {
	AuditID     int       `json:"audit_id"`
	StatusInt   int       `json:"status_int"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	model.PageSize
	model.DateRequest
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