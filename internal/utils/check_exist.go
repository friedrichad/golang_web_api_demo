package utils

import "gorm.io/gorm"

func Exists(db *gorm.DB, query string, args ...interface{}) (bool, error) {
	var exists bool

	err := db.Raw(
		"SELECT EXISTS(" + query + ")",
		args...,
	).Scan(&exists).Error

	return exists, err
}