package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type CustomerRequest struct {
	CustomerID   *int       `json:"customer_id"`
	CustomerName *string    `json:"customer_name"`
	Phone        *string    `json:"phone"`
	Email        *string    `json:"email"`
	Address      *string    `json:"address"`
	StatusInt    *int       `json:"status_int"`
	model.PageSize
	model.DateRequest
}
type CustomerResponse struct {
	CustomerID   int       `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	StatusInt    int       `json:"status_int"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    int       `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}