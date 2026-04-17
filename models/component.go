package models

type Component struct {
	ComponentID  int     `gorm:"primaryKey;autoIncrement" json:"component_id"`
	MetadataJSON string  `gorm:"type:json;not null" json:"metadata_json"`
	Unit         string  `gorm:"size:50;not null" json:"unit"`
	UnitPrice    float64 `gorm:"not null" json:"unit_price"`
}
