package dtos

import(
	"time"
)

type UserResp struct {
	UserUUID     string    `gorm:"column:user_uuid;primaryKey" json:"user_uuid"`
	Username     string    `gorm:"column:username" json:"username"`
	DisplayName  string    `gorm:"column:display_name" json:"display_name"`
	Email        string    `gorm:"column:email" json:"email"`
	StatusInt    int32     `gorm:"column:status_int" json:"status_int"`
	CreatedBy    int32     `gorm:"column:created_by" json:"created_by"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy    int32     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}
type UserReq struct {
	UserUUID     string    `gorm:"column:user_uuid;primaryKey" json:"user_uuid"`
	Username     string    `gorm:"column:username" json:"username"`
	DisplayName  string    `gorm:"column:display_name" json:"display_name"`
	Email        string    `gorm:"column:email" json:"email"`
	StatusInt    int32     `gorm:"column:status_int" json:"status_int"`
	CreatedBy    int32     `gorm:"column:created_by" json:"created_by"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy    int32     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}