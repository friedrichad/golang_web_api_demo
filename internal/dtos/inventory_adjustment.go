package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// InventoryAdjustmentFilter - GET request with query parameters
type InventoryAdjustmentFilter struct {
	AdjustmentID *int `form:"adjustment_id"`
	AuditID      *int `form:"audit_id"`
	StatusInt    *int `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// InventoryAdjustmentCreate - POST request body
type InventoryAdjustmentCreate struct {
	AuditID     int                               `json:"audit_id" binding:"required"`
	Description string                            `json:"description"`
	Note        string                            `json:"note"`
	Details     []InventoryAdjustmentDetailCreate `json:"details"`
}

// InventoryAdjustmentUpdate - PUT request body
type InventoryAdjustmentUpdate struct {
	AdjustmentID int                               `json:"adjustment_id" binding:"required"`
	ApprovedID   int                               `json:"approved_id"`
	Description  string                            `json:"description"`
	StatusInt    int                               `json:"status_int"`
	Note         string                            `json:"note"`
	UpdatedBy    int                               `json:"updated_by"`
	Details      []InventoryAdjustmentDetailUpdate `json:"details"`
}

type InventoryAdjustmentResponse struct {
	AdjustmentID int                                 `json:"adjustment_id"`
	AuditID      int                                 `json:"audit_id"`
	ApprovedID   int                                 `json:"approved_id"`
	Description  string                              `json:"description"`
	StatusInt    int                                 `json:"status_int"`
	Note         string                              `json:"note"`
	CreatedBy    int                                 `json:"created_by"`
	CreatedAt    time.Time                           `json:"created_at"`
	UpdatedBy    int                                 `json:"updated_by"`
	UpdatedAt    time.Time                           `json:"updated_at"`
	Details      []InventoryAdjustmentDetailResponse `json:"details"`
}
