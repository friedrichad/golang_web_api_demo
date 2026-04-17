package models

import "time"

type TransferRequest struct {
	RequestID       int        `gorm:"primaryKey;autoIncrement" json:"request_id"`
	Description     string     `gorm:"type:text;not null" json:"description"`
	TypeInt         int        `gorm:"not null" json:"type_int"`
	CreationTime    time.Time  `gorm:"not null" json:"creation_time"`
	ExecutionTime   *time.Time `json:"execution_time"`
	StatusInt       int        `gorm:"not null" json:"status_int"`
	CreatorID       int        `gorm:"not null" json:"creator_id"`
	ApproverID      *int       `json:"approver_id"`
	WarehouseFromID *int       `json:"warehouse_from_id"`
	WarehouseToID   *int       `json:"warehouse_to_id"`
	CustomerID      *int       `json:"customer_id"`
}
