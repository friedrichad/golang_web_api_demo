package models

type User struct {
	UserID       int    `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username     string `gorm:"size:25;not null" json:"username"`
	DisplayName  string `gorm:"size:100;not null" json:"display_name"`
	Email        string `gorm:"size:320;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"` // Không trả về password trong JSON
	StatusInt    int    `gorm:"not null" json:"status_int"`
}
