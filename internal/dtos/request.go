package dtos

import (
	"time"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type RequestRequest struct{
	RequestID     int       `json:"request_id"`
	RequestType   string    `json:"request_type"`
	StatusInt     int       `json:"status_int"`
	model.PageSize
	model.DateRequest
}
type RequestReponse struct{
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
	UpdatedBy     int   	`json:"updated_by"`
}
type RequestUpdate struct{
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
}