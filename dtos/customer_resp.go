package dtos

import("time")

type Customer struct {
	CustomerID   int32    `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	TypeInt      int32     `json:"type_int"`
}