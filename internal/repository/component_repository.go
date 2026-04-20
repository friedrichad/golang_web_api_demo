package repository

import (
	"fmt"

	dtos "github.com/friedrichad/golang_web_api_demo/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/models"
	"gorm.io/gorm"
)

type IComponentRepository interface {
	GetComponentByID(componentID int32) (*dtos.ComponentResponse, error)
	GetComponents() ([]*dtos.ComponentResponse, error)
	CreateComponent(component *dtos.ComponentRequest) (*dtos.ComponentResponse, error)
}

type ComponentRepository struct {
	DB *gorm.DB
}

func (s *ComponentRepository) GetComponentByID(componentId int32) (*dtos.ComponentResponse, error) {
	var component models.Component
	result := s.DB.Where("component_id = ?", componentId).First(&component)
	if result.Error != nil {
		return nil, result.Error
	}
	return &dtos.ComponentResponse{
		ComponentID:  component.ComponentID,
		MetadataJSON: component.MetadataJSON,
		Unit:         component.Unit,
		UnitPrice:    component.UnitPrice,
	}, nil
}
func (s *ComponentRepository) GetComponents() ([]*dtos.ComponentResponse, error) {
	var components []models.Component
	result := s.DB.Find(&components)
	if result.Error != nil {
		return nil, result.Error
	}
	var componentResps []*dtos.ComponentResponse
	for _, component := range components {
		componentResps = append(componentResps, &dtos.ComponentResponse{
			ComponentID:  component.ComponentID,
			MetadataJSON: component.MetadataJSON,
			Unit:         component.Unit,
			UnitPrice:    component.UnitPrice,
		})
	}
	return componentResps, nil
}
func (s *ComponentRepository) CreateComponent(componentReq *dtos.ComponentRequest) (*dtos.ComponentResponse, error) {
	var bin models.Bin
	result := s.DB.Where("bin_id =?", componentReq.BinID).First(&bin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bin with ID %d not found", componentReq.BinID)
		}
		return nil, result.Error
	}
	var category models.Componentcategory
	result = s.DB.Where("category_id =?", componentReq.ComponentCategoryID).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("component category with ID %d not found", componentReq.ComponentCategoryID)
		}
		return nil, result.Error
	}
	component := models.Component{
		MetadataJSON: componentReq.MetadataJSON,
		Unit:         componentReq.Unit,
		UnitPrice:    componentReq.UnitPrice,
	}
	result = s.DB.Create(&component)
	if result.Error != nil {
		return nil, result.Error
	}
	componentBin := models.Componentbin{
		ComponentID: component.ComponentID,
		Quantity:    componentReq.Quantity,
		BinID:       componentReq.BinID,
	}
	result = s.DB.Create(&componentBin)
	if result.Error != nil {
		return nil, result.Error
	}
	componentCategory := models.CC{
		ComponentID: component.ComponentID,
		CategoryID:  category.CategoryID,
	}
	result = s.DB.Create(&componentCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &dtos.ComponentResponse{
		ComponentID:  component.ComponentID,
		MetadataJSON: component.MetadataJSON,
		Unit:         component.Unit,
		UnitPrice:    component.UnitPrice,
		AddComponentToBinRequest: dtos.AddComponentToBinRequest{
			BinID:    componentBin.BinID,
			Quantity: componentBin.Quantity,
		},
		ComponentCategoryResponse: dtos.ComponentCategoryResponse{
			CategoryID:   category.CategoryID,
			CategoryName: category.CategoryName,
		},
	}, nil
}
