package models

type Bin struct {
	BinID               int       `gorm:"primaryKey;autoIncrement" json:"bin_id"`
	LocationInWarehouse string    `gorm:"type:text;not null" json:"location_in_warehouse"`
	StatusInt           int       `gorm:"not null" json:"status_int"`
	WarehouseID         int       `gorm:"not null" json:"warehouse_id"`
	Warehouse           Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
}
