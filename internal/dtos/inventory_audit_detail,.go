package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type InventoryAuditDetailResponse struct {
	AuditDetailID      int       `json:"audit_detail_id"`
	AuditID            int       `json:"audit_id"`
	ComponentID        int       `json:"component_id"`
	BinID              int       `json:"bin_id"`
	SystemQuantity     float64   `json:"system_quantity"`
	ActualQuantity     float64   `json:"actual_quantity"`
	DifferenceQuantity float64   `json:"difference_quantity"`
	Note               string    `json:"note"`
	CreatedBy          int       `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
}

type InventoryAuditDetailRequest struct {
	AuditDetailID      int       `json:"audit_detail_id"`
	AuditID            int       `json:"audit_id"`
	model.PageSize
	model.DateRequest
}
type InventoryAuditDetailUpdate struct {
	AuditDetailID      int       `json:"audit_detail_id"`
	AuditID            int       `json:"audit_id"`
	ComponentID        int       `json:"component_id"`
	BinID              int       `json:"bin_id"`
	SystemQuantity     float64   `json:"system_quantity"`
	ActualQuantity     float64   `json:"actual_quantity"`
	DifferenceQuantity float64   `json:"difference_quantity"`
	Note               string    `json:"note"`
	CreatedBy          int       `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
}