package models

type ComponentBin struct {
	ComponentBinID int       `gorm:"primaryKey;autoIncrement" json:"component_bin_id"`
	Quantity       float64   `gorm:"not null" json:"quantity"`
	ComponentID    int       `gorm:"not null" json:"component_id"`
	BinID          int       `gorm:"not null" json:"bin_id"`
	Component      Component `gorm:"foreignKey:ComponentID"`
	Bin            Bin       `gorm:"foreignKey:BinID"`
}