package dtos

import (
	"time"
	"github.com/friedrichad/golang_web_api_demo/internal/model")

type InventoryAdjustmentRequest struct{
	AdjustmentID int       `json:"adjustment_id"`
	AuditID      int       `json:"audit_id"`
	StatusInt    int       `json:"status_int"`
	model.PageSize
	model.DateRequest
}
type InventoryAdjustmentResponse struct{
	AdjustmentID int       `json:"adjustment_id"`
	AuditID      int       `json:"audit_id"`
	ApprovedID   int       `json:"approved_id"`
	Description  string    `json:"description"`
	ApprovedTime time.Time `json:"approved_time"`
	StatusInt    int       `json:"status_int"`
	Note         string    `json:"note"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}