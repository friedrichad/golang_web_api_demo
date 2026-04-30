package dtos

import "github.com/friedrichad/golang_web_api_demo/internal/model"

type RequestDetailDTO struct {
	RequestDetailID int     `json:"request_detail_id"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
}

type RequestDetailRequest struct {
	RequestDetailID int     `json:"request_detail_id"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	model.PageSize
	model.DateRequest
}
type RequestDetailResponse struct {
	RequestDetailID int     `json:"request_detail_id"`
	RequestID       int     `json:"request_id"`
	ComponentID     int     `json:"component_id"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
}