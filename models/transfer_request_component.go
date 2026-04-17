package models

type TransferRequestComponent struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID   int     `gorm:"not null" json:"request_id"`
	ComponentID int     `gorm:"not null" json:"component_id"`
	Quantity    float64 `gorm:"not null" json:"quantity"`
	UnitPrice   float64 `gorm:"not null" json:"unit_price"`
}