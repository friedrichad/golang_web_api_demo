package models

type Warehouse struct {
	WarehouseID      int    `gorm:"primaryKey;autoIncrement" json:"warehouse_id"`
	WarehouseName    string `gorm:"size:50;not null" json:"warehouse_name"`
	Description      string `gorm:"type:text;not null" json:"description"`
	PhysicalLocation string `gorm:"type:text;not null" json:"physical_location"`
	ImageURL         string `gorm:"size:1000" json:"image_url"`
}