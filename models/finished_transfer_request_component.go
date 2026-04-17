package models

type FinishedTransferRequestComponent struct {
	FinishedTransferRequestComponentID int     `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID                          int     `gorm:"not null" json:"request_id"`
	ComponentID                        int     `gorm:"not null" json:"component_id"`
	BinID                              int     `gorm:"not null" json:"bin_id"`
	Quantity                           float64 `gorm:"not null" json:"quantity"`
	TypeInt                            int     `gorm:"not null" json:"type_int"`
}
