package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
)

// RequestFilter - GET request with query parameters
type RequestFilter struct {
	RequestID   *int    `form:"request_id"`
	RequestType *string `form:"request_type"`
	StatusInt   *int    `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// RequestCreate - POST request body (creates Request header only)
// RequestDetail must be created separately via CreateRequestDetail
type RequestCreate struct {
	RequestType *int      `json:"request_type" binding:"required"`
	Description string    `json:"description"`
	WarehouseID int       `json:"warehouse_id" binding:"required"`
	PartnerID   int       `json:"partner_id"`
	Note        string    `json:"note"`
	ExpiredDate time.Time `json:"expired_date"`
}

// RequestUpdate - PUT request body
type RequestUpdate struct {
	RequestID     int    `json:"request_id" binding:"required"`
	RequestType   *int   `json:"request_type"`
	Description   string `json:"description"`
	WarehouseID   int    `json:"warehouse_id"`
	PerformedByID int    `json:"performed_by_id"`
	ApproverID    int    `json:"approver_id"`
	PartnerID     int    `json:"partner_id"`
	StatusInt     int    `json:"status_int"`
	Note          string `json:"note"`
	UpdatedBy     int    `json:"updated_by"`
}

// RequestResponse - Response body
type RequestResponse struct {
	RequestID     int       `json:"request_id"`
	RequestType   int       `json:"request_type"`
	Description   string    `json:"description"`
	WarehouseID   int       `json:"warehouse_id"`
	PerformedByID int       `json:"performed_by_id"`
	ApproverID    int       `json:"approver_id"`
	PartnerID     int       `json:"partner_id"`
	ExpiredDate   time.Time `json:"expired_date"`
	StatusInt     int       `json:"status_int"`
	Note          string    `json:"note"`
	Reason        string    `json:"reason"`
	CreatedBy     int       `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedBy     int       `json:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type ApprovalRequest struct {
	RequestID  int    `json:"request_id"`
	ApproverID int    `json:"approver_id"`
	StatusInt  int    `json:"status_int"`
	Reason     string `json:"reason"`
}
type ConfirmRequest struct {
	RequestID int `json:"request_id"`
	StatusInt int `json:"status_int"`
}

// Verify validates the RequestCreate struct (header validation only).
func (r *RequestCreate) Verify() error {
	if r.RequestType == nil {
		return fmt.Errorf("RequestType is required")
	}
	if !constants.IsValidWarehouseRequestType(*r.RequestType) {
		return fmt.Errorf("Invalid RequestType. Valid types: %v", *r.RequestType)
	}
	if r.WarehouseID == 0 {
		return fmt.Errorf("WarehouseID is required")
	}
	// RequestDetail validation moved to CreateRequestDetail endpoint
	return nil
}

// Verify validates the RequestUpdate struct.
func (r *RequestUpdate) Verify() error {
	if r.RequestID == 0 {
		return fmt.Errorf("RequestID is required")
	}
	if r.RequestType != nil && !constants.IsValidWarehouseRequestType(*r.RequestType) {
		return fmt.Errorf("Invalid RequestType. Valid types: %v", *r.RequestType)
	}
	return nil
}
