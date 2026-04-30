package dtos

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type ComponentRequest struct {
	ComponentID       int                    `json:"component_id"`
	ComponentName     string                 `json:"component_name"`
	MetadataJSON      string                 `json:"metadata_json"`
	Unit              string                 `json:"unit"`
	UnitPrice         float64                `json:"unit_price"`
	ComponentCategory []ComponentCategoryDTO `json:"component_category"`
	CreatedBy         int                    `json:"created_by"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedBy         int                    `json:"updated_by"`
	UpdatedAt         time.Time              `json:"updated_at"`
	model.PageSize
	model.DateRequest
}

type ComponentResponse struct {
	ComponentID       int                    `json:"component_id"`
	ComponentName     string                 `json:"component_name"`
	MetadataJSON      string                 `json:"metadata_json"`
	Unit              string                 `json:"unit"`
	UnitPrice         float64                `json:"unit_price"`
	Quantity          float64                `json:"quantity"`
	ComponentCategory []ComponentCategoryDTO `json:"component_category"`
	CreatedBy         int                    `json:"created_by"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedBy         int                    `json:"updated_by"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type ComponentCategoryDTO struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}
type ComponentMetadata struct {
	ImageUrl    string `json:"image_url"`
	Description string `json:"description"`
}
type UploadFileInfoDTO struct {
	FileName    string    `json:"file_name"`
	ObjectName  string    `json:"object_name"`
	Bucket      string    `json:"bucket"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	ETag        string    `json:"etag"`
	UploadedAt  time.Time `json:"uploaded_at"`
	URL         string    `json:"url"`
}
