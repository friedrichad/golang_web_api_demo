package shared

type UserPermissionScope struct {
	Scope       string `gorm:"column:scope" json:"scope"`
	ExpiredDate int64  `gorm:"column:expired_date" json:"expired_date"`
}

type UserInfo struct{
	UserId int `json:"user_id"`
	UserName string `json:"user_name"`
	PositionInfo PositionInfo `json:"position_info"`
	IsOP int `json:"is_op"`
}

type PositionInfo struct{
	PositionId int `json:"position_id"`
	PositionName string `json:"position_name"`
	PositionLevel int `json:"position_level"`
}