package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// RequestFilter - GET request with query parameters
type RequestFilter struct {
	RequestID   int    `form:"request_id"`
	RequestType string `form:"request_type"`
	StatusInt   int    `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// RequestCreate - POST request body
type RequestCreate struct {
	RequestType   string    `json:"request_type" binding:"required"`
	Description   string    `json:"description"`
	WarehouseID   int       `json:"warehouse_id" binding:"required"`
	BinFrom       int       `json:"bin_from"`
	BinTo         int       `json:"bin_to"`
	PerformedByID int       `json:"performed_by_id"`
	PartnerID     int       `json:"partner_id"`
	RequestDate   time.Time `json:"request_date"`
	Note          string    `json:"note"`
}

// RequestUpdate - PUT request body
type RequestUpdate struct {
	RequestID     int       `json:"request_id" binding:"required"`
	RequestType   string    `json:"request_type"`
	Description   string    `json:"description"`
	WarehouseID   int       `json:"warehouse_id"`
	BinTo         int       `json:"bin_to"`
	BinFrom       int       `json:"bin_from"`
	PerformedByID int       `json:"performed_by_id"`
	ApproverID    int       `json:"approver_id"`
	PartnerID     int       `json:"partner_id"`
	RequestDate   time.Time `json:"request_date"`
	StatusInt     int       `json:"status_int"`
	Note          string    `json:"note"`
	UpdatedBy     int       `json:"updated_by"`
}

// RequestResponse - Response body
type RequestResponse struct {
	RequestID     int       `json:"request_id"`
	RequestType   string    `json:"request_type"`
	Description   string    `json:"description"`
	WarehouseID   int       `json:"warehouse_id"`
	BinTo         int       `json:"bin_to"`
	BinFrom       int       `json:"bin_from"`
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
