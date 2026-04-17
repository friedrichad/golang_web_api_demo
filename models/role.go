package models
type Role struct {
	RoleID      int    `gorm:"primaryKey;autoIncrement" json:"role_id"`
	RoleName    string `gorm:"size:50;not null" json:"role_name"`
	Description string `gorm:"type:text" json:"description"`
}