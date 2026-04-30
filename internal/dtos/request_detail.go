package dtos

import "github.com/friedrichad/golang_web_api_demo/internal/model"

type RequestDetailDTO struct {
	RequestDetailID int     `json:"request_detail_id"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
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
}

// RequestDetailUpdate - PUT request body
type RequestDetailUpdate struct {
	RequestDetailID int     `json:"request_detail_id" binding:"required"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	UpdatedBy       int     `json:"updated_by"`
}

type RequestDetailResponse struct {
	RequestDetailID int     `json:"request_detail_id"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
}
