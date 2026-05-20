package shared

type UserPermissionScope struct {
	Scope       string `gorm:"column:scope" json:"scope"`
	ExpiredDate int64  `gorm:"column:expired_date" json:"expired_date"`
}
