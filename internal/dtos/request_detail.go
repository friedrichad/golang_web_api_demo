package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type RequestDetailDTO struct {
	RequestID   *int     `json:"request_id"`
	ComponentID *int     `json:"component_id"`
	Quantity    *int     `json:"quantity"`
	UnitPrice   *float64 `json:"unit_price"`
	BinFromID   *int     `json:"bin_from_id"`
	BinToID     *int     `json:"bin_to_id"`
}

// RequestDetailFilter - GET request with query parameters
type RequestDetailFilter struct {
	RequestDetailID int `form:"request_detail_id"`
	RequestID       int `form:"request_id"`
	ComponentID     int `form:"component_id"`
	model.PageSize
	model.DateRequest
}

// RequestDetailCreate - POST request body
type RequestDetailCreate struct {
	RequestID   int     `json:"request_id" binding:"required"`
	ComponentID int     `json:"component_id" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	BinFromID   int     `json:"bin_from_id"`
	BinToID     int     `json:"bin_to_id"`
}

// Verify validates the RequestDetailCreate struct.
func (r *RequestDetailCreate) Verify() error {
	if r.RequestID == 0 {
		return fmt.Errorf("RequestID is required")
	}
	if r.ComponentID == 0 {
		return fmt.Errorf("ComponentID is required")
	}
	if r.Quantity <= 0 {
		return fmt.Errorf("Quantity must be greater than 0")
	}
	return nil
}

// RequestDetailUpdate - PUT request body
type RequestDetailUpdate struct {
	RequestDetailID int     `json:"request_detail_id" binding:"required"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	BinFromID       int     `json:"bin_from_id"`
	BinToID         int     `json:"bin_to_id"`
	UpdatedBy       int     `json:"updated_by"`
}

// Verify validates the RequestDetailUpdate struct.
func (r *RequestDetailUpdate) Verify() error {
	if r.RequestDetailID == 0 {
		return fmt.Errorf("RequestDetailID is required")
	}
	return nil
}

type RequestDetailResponse struct {
	RequestDetailID int       `json:"request_detail_id"`
	RequestID       int       `json:"request_id"`
	ComponentID     int       `json:"component_id"`
	Quantity        int       `json:"quantity"`
	UnitPrice       float64   `json:"unit_price"`
	BinFromID       int       `json:"bin_from_id"`
	BinToID         int       `json:"bin_to_id"`
	CreatedBy       int       `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedBy       int       `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}
