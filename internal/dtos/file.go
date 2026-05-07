package dtos

type UploadBase64Request struct {
	FileName string `json:"file_name" binding:"required"`
	Data     string `json:"data" binding:"required"`
}
type UploadResponse struct {
	FileName string `json:"file_name"`
	FileUrl  string `json:"file_url"`
	Size     int64  `json:"size"`
}