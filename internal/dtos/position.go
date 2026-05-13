package dtos

import "github.com/friedrichad/golang_web_api_demo/internal/model"

type PositionFilter struct {
	PositionID   int    `json:"position_id"`
	PositionName string `json:"position_name"`
	Description  string `json:"description"`
	model.PageSize
	model.DateRequest
}
type PositionCreate struct {
	PositionName string `json:"position_name" binding:"required"`
	Description  string `json:"description"`
}

type PositionUpdate struct {
	PositionID   int    `json:"position_id" binding:"required"`
	PositionName string `json:"position_name"`
	Description  string `json:"description"`
}

type PositionResponse struct {
	PositionID   int    `json:"position_id"`
	PositionName string `json:"position_name"`
	Description  string `json:"description"`
	CreatedBy    int    `json:"created_by"`
	CreatedAt    string `json:"created_at"`
	UpdatedBy    int    `json:"updated_by"`
	UpdatedAt    string `json:"updated_at"`
}
