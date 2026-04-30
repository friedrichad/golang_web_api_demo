package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type ComponentBinCreate struct {
	Quantity    float64 `json:"quantity"`
	ComponentID int32   `json:"component_id"`
	BinID       int32   `json:"bin_id"`
}

type ComponentBinUpdate struct {
	ComponentBinID int32   `json:"component_bin_id"`
	Quantity       float64 `json:"quantity"`
	ComponentID    int32   `json:"component_id"`
	BinID          int32   `json:"bin_id"`
}

type ComponetBinFilter struct {
	Quantity    float64 `form:"quantity"`
	ComponentID int32   `form:"component_id"`
	BinID       int32   `form:"bin_id"`
	model.PageSize
	model.DateRequest
}
type ComponentBinResponse struct {
	ComponentBinID int32     `json:"component_bin_id"`
	Quantity       float64   `json:"quantity"`
	ComponentID    int32     `json:"component_id"`
	BinID          int32     `json:"bin_id"`
	CreatedBy      int32     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedBy      int32     `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}
