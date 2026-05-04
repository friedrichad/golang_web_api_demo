package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type ComponentBinCreate struct {
	Quantity    float64 `json:"quantity"`
	ComponentID int     `json:"component_id"`
	BinID       int     `json:"bin_id"`
}

type ComponentBinUpdate struct {
	ComponentBinID int     `json:"component_bin_id"`
	Quantity       float64 `json:"quantity"`
	ComponentID    int     `json:"component_id"`
	BinID          int     `json:"bin_id"`
}

type ComponetBinFilter struct {
	Quantity    *float64 `form:"quantity"`
	ComponentID *int     `form:"component_id"`
	BinID       *int     `form:"bin_id"`
	model.PageSize
	model.DateRequest
}
type ComponentBinResponse struct {
	ComponentBinID int       `json:"component_bin_id"`
	Quantity       float64   `json:"quantity"`
	ComponentID    int       `json:"component_id"`
	BinID          int       `json:"bin_id"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedBy      int       `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}
