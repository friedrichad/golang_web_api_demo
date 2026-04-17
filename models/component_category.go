package models

type ComponentCategory struct {
	CategoryID   int    `gorm:"primaryKey;autoIncrement" json:"category_id"`
	CategoryName string `gorm:"size:50;not null" json:"category_name"`
}