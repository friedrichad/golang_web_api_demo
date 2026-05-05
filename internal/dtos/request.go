package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// RequestFilter - GET request with query parameters
type RequestFilter struct {
	RequestID   *int    `form:"request_id"`
	RequestType *string `form:"request_type"`
	StatusInt   *int    `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// RequestCreate - POST request body
type RequestCreate struct {
	RequestType   string             `json:"request_type" binding:"required"`
	Description   string             `json:"description"`
	WarehouseID   int                `json:"warehouse_id" binding:"required"`
	PerformedByID int                `json:"performed_by_id"`
	PartnerID     int                `json:"partner_id"`
	RequestDate   time.Time          `json:"request_date"`
	Note          string             `json:"note"`
	RequestDetail []RequestDetailDTO `json:"request_detail" binding:"required"`
}

// RequestUpdate - PUT request body
type RequestUpdate struct {
	RequestID     int    `json:"request_id" binding:"required"`
	RequestType   string `json:"request_type"`
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
	RequestType   string    `json:"request_type"`
	Description   string    `json:"description"`
	WarehouseID   int       `json:"warehouse_id"`
	PerformedByID int       `json:"performed_by_id"`
	ApproverID    int       `json:"approver_id"`
	PartnerID     int       `json:"partner_id"`
	RequestDate   time.Time `json:"request_date"`
	StatusInt     int       `json:"status_int"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
	CreateBy      int       `json:"create_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     int       `json:"updated_by"`
}
type ApprovalRequest struct {
    RequestID    int       `json:"request_id"`
    ApproverID   int       `json:"approver_id"`
    StatusInt    int       `json:"status_int"`
    Note         string    `json:"note"`
}
type ConfirmRequest struct {
	RequestID    int       `json:"request_id"`
	StatusInt    int       `json:"status_int"`
}

// Verify validates the RequestCreate struct.
func (r *RequestCreate) Verify() error {
	if r.RequestType == "" {
		return fmt.Errorf("RequestType is required")
	}
	if r.WarehouseID == 0 {
		return fmt.Errorf("WarehouseID is required")
	}
	if len(r.RequestDetail) == 0 {
		return fmt.Errorf("At least one RequestDetail is required")
	}
	return nil
}

// Verify validates the RequestUpdate struct.
func (r *RequestUpdate) Verify() error {
	if r.RequestID == 0 {
		return fmt.Errorf("RequestID is required")
	}
	return nil
}

