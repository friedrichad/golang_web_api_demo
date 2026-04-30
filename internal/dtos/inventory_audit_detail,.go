package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// InventoryAuditDetailFilter - GET request with query parameters
type InventoryAuditDetailFilter struct {
	AuditDetailID int `form:"audit_detail_id"`
	AuditID       int `form:"audit_id"`
	ComponentID   int `form:"component_id"`
	BinID         int `form:"bin_id"`
	model.PageSize
	model.DateRequest
}

// InventoryAuditDetailCreate - POST request body
type InventoryAuditDetailCreate struct {
	AuditID            int     `json:"audit_id" binding:"required"`
	ComponentID        int     `json:"component_id" binding:"required"`
	BinID              int     `json:"bin_id" binding:"required"`
	SystemQuantity     float64 `json:"system_quantity" binding:"required"`
	ActualQuantity     float64 `json:"actual_quantity" binding:"required"`
	DifferenceQuantity float64 `json:"difference_quantity"`
	Note               string  `json:"note"`
}

// InventoryAuditDetailUpdate - PUT request body
type InventoryAuditDetailUpdate struct {
	AuditDetailID      int     `json:"audit_detail_id" binding:"required"`
	AuditID            int     `json:"audit_id"`
	ComponentID        int     `json:"component_id"`
	BinID              int     `json:"bin_id"`
	SystemQuantity     float64 `json:"system_quantity"`
	ActualQuantity     float64 `json:"actual_quantity"`
	DifferenceQuantity float64 `json:"difference_quantity"`
	Note               string  `json:"note"`
	UpdatedBy          int     `json:"updated_by"`
}

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
