package models

import "time"

type Customer struct {
	CustomerID   int       `gorm:"primaryKey;autoIncrement" json:"customer_id"`
	CustomerName string    `gorm:"size:100;not null" json:"customer_name"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Email        string    `gorm:"size:320;unique" json:"email"`
	Address      string    `gorm:"type:text" json:"address"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	TypeInt      int       `json:"type_int"`
}
