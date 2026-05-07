package constants

import "errors"

const maxFileSize = 5 * 1024 * 1024

var allowedTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}

func ValidateFileSize(size int64) error {

	if size > maxFileSize {
		return errors.New("file too large")
	}

	return nil
}

func ValidateFileType(contentType string) error {

	if !allowedTypes[contentType] {
		return errors.New("file type not allowed")
	}

	return nil
}
