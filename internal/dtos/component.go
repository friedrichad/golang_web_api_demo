package dtos

import (
	"fmt"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

// ComponentFilter - GET request with query parameters
type ComponentFilter struct {
	ComponentID   *int    `form:"component_id"`
	ComponentName *string `form:"component_name"`
	model.PageSize
	model.DateRequest
}

// ComponentCreate - POST request body
type ComponentCreate struct {
	ComponentName     string                 `json:"component_name" binding:"required"`
	MetadataJSON      string                 `json:"metadata_json"`
	Unit              string                 `json:"unit" binding:"required"`
	UnitPrice         float64                `json:"unit_price" binding:"required"`
	ComponentCategory []ComponentCategoryDTO `json:"component_category"`
}

// Verify validates the ComponentCreate struct.
func (c *ComponentCreate) Verify() error {
	if c.ComponentName == "" {
		return fmt.Errorf("ComponentName is required")
	}
	if c.Unit == "" {
		return fmt.Errorf("Unit is required")
	}
	if c.UnitPrice <= 0 {
		return fmt.Errorf("UnitPrice must be greater than 0")
	}
	return nil
}

// ComponentUpdate - PUT request body
type ComponentUpdate struct {
	ComponentID       int                    `json:"component_id" binding:"required"`
	ComponentName     string                 `json:"component_name"`
	MetadataJSON      string                 `json:"metadata_json"`
	Unit              string                 `json:"unit"`
	UnitPrice         float64                `json:"unit_price"`
	ComponentCategory []ComponentCategoryDTO `json:"component_category"`
	UpdatedBy         int                    `json:"updated_by"`
}

// Verify validates the ComponentUpdate struct.
func (c *ComponentUpdate) Verify() error {
	if c.ComponentID == 0 {
		return fmt.Errorf("ComponentID is required")
	}
	return nil
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
