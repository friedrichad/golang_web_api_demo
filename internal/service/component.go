package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IComponentService interface {
	GetAllComponents(c *gin.Context) ([]dtos.ComponentResponse, int, *common.Error)
	GetComponentById(c *gin.Context) (*dtos.ComponentResponse, *common.Error)
	CreateComponent(c *gin.Context) (*dtos.ComponentResponse, *common.Error)
	UpdateComponent(c *gin.Context) *common.Error
	DeleteComponent(c *gin.Context) *common.Error
}

type ComponentService struct {
	componentRepo repository.IComponentRepository
	db            *gorm.DB
}

var componentService IComponentService

func NewComponentService() IComponentService {
	if componentService == nil {
		componentService = &ComponentService{
			componentRepo: repository.NewComponentRepository(),
			db:            db.Instance,
		}
	}
	return componentService
}

func (s *ComponentService) GetAllComponents(c *gin.Context) ([]dtos.ComponentResponse, int, *common.Error) {
	var query dtos.ComponentFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}

	components, total, err := s.componentRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	res := make([]dtos.ComponentResponse, len(components))
	for i, comp := range components {
		resp := modelToComponentResponse(&comp)
		resp.ComponentCategory = s.getComponentCategories(int(comp.ComponentID))
		// Quantity can be joined from inventory, but for simplicty we set 0 if not handled by repo
		res[i] = resp
	}

	return res, total, nil
}

func (s *ComponentService) GetComponentById(c *gin.Context) (*dtos.ComponentResponse, *common.Error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	comp, err := s.componentRepo.GetById(id)
	if err != nil || comp == nil {
		return nil, common.NotFound
	}

	res := modelToComponentResponse(comp)
	res.ComponentCategory = s.getComponentCategories(id)
	return &res, nil
}

func (s *ComponentService) CreateComponent(c *gin.Context) (*dtos.ComponentResponse, *common.Error) {
	var req dtos.ComponentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	comp := &model.Component{
		ComponentName: req.ComponentName,
		MetadataJSON:  req.MetadataJSON,
		Unit:          req.Unit,
		UnitPrice:     req.UnitPrice,
		CreatedAt:     time.Now(),
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comp).Error; err != nil {
			return err
		}

		for _, cat := range req.ComponentCategory {
			mapEntry := &model.ComponentCategoryMap{
				ComponentID: comp.ComponentID,
				CategoryID:  int32(cat.CategoryID),
			}
			if err := tx.Create(mapEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, common.SystemError
	}

	res := modelToComponentResponse(comp)
	res.ComponentCategory = s.getComponentCategories(int(comp.ComponentID))
	return &res, nil
}

func (s *ComponentService) UpdateComponent(c *gin.Context) *common.Error {
	var req dtos.ComponentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	comp, err := s.componentRepo.GetById(req.ComponentID)
	if err != nil || comp == nil {
		return common.NotFound
	}

	if req.ComponentName != "" {
		comp.ComponentName = req.ComponentName
	}
	if req.MetadataJSON != "" {
		comp.MetadataJSON = req.MetadataJSON
	}
	if req.Unit != "" {
		comp.Unit = req.Unit
	}
	if req.UnitPrice != 0 {
		comp.UnitPrice = req.UnitPrice
	}
	if req.UpdatedBy != 0 {
		comp.UpdatedBy = int32(req.UpdatedBy)
	}
	comp.UpdatedAt = time.Now()

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(comp).Error; err != nil {
			return err
		}

		// Delete existing associations
		if err := tx.Where("component_id = ?", comp.ComponentID).Delete(&model.ComponentCategoryMap{}).Error; err != nil {
			return err
		}

		// Insert new associations
		for _, cat := range req.ComponentCategory {
			mapEntry := &model.ComponentCategoryMap{
				ComponentID: comp.ComponentID,
				CategoryID:  int32(cat.CategoryID),
			}
			if err := tx.Create(mapEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *ComponentService) DeleteComponent(c *gin.Context) *common.Error {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("component_id IN ?", ids).Delete(&model.ComponentCategoryMap{}).Error; err != nil {
			return err
		}
		if err := tx.Where("component_id IN ?", ids).Delete(&model.Component{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *ComponentService) getComponentCategories(componentId int) []dtos.ComponentCategoryDTO {
	var categories []dtos.ComponentCategoryDTO
	s.db.Model(&model.ComponentCategory{}).
		Select("component_category.category_id, component_category.category_name").
		Joins("INNER JOIN component_category_map ON component_category.category_id = component_category_map.category_id").
		Where("component_category_map.component_id = ?", componentId).
		Scan(&categories)
	return categories
}

func modelToComponentResponse(c *model.Component) dtos.ComponentResponse {
	return dtos.ComponentResponse{
		ComponentID:   int(c.ComponentID),
		ComponentName: c.ComponentName,
		MetadataJSON:  c.MetadataJSON,
		Unit:          c.Unit,
		UnitPrice:     c.UnitPrice,
		CreatedBy:     int(c.CreatedBy),
		CreatedAt:     c.CreatedAt,
		UpdatedBy:     int(c.UpdatedBy),
		UpdatedAt:     c.UpdatedAt,
	}
}