package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// CustomerFilter - GET request with query parameters
type CustomerFilter struct {
	CustomerName *string `form:"customer_name"`
	Phone        *string `form:"phone"`
	Email        *string `form:"email"`
	StatusInt    *int    `form:"status_int"`
	model.PageSize
	model.DateRequest
}

// CustomerCreate - POST request body
type CustomerCreate struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email" binding:"email"`
	Address      string `json:"address"`
}

// CustomerCreateVerify - Verify method for CustomerCreate
func (c *CustomerCreate) Verify() error {
	if c.CustomerName == "" {
		return fmt.Errorf("CustomerName is required")
	}
	if c.Email == "" {
		return fmt.Errorf("Email is required")
	}
	return nil
}

// CustomerUpdate - PUT request body
type CustomerUpdate struct {
	CustomerID   int    `json:"customer_id" binding:"required"`
	CustomerName string `json:"customer_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	StatusInt    *int   `json:"status_int"`
	UpdatedBy    int    `json:"updated_by"`
}

// CustomerUpdateVerify - Verify method for CustomerUpdate
func (c *CustomerUpdate) Verify() error {
	if c.CustomerID == 0 {
		return fmt.Errorf("CustomerID is required")
	}
	return nil
}

// CustomerResponse - Response body
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
