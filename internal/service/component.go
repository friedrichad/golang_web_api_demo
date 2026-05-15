package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IComponentService interface {
	GetAllComponents(c *gin.Context) ([]model.ComponentResponse, int, *common.Error)
	GetComponentById(c *gin.Context) (*model.ComponentResponse, *common.Error)
	CreateComponent(c *gin.Context) (*model.ComponentResponse, *common.Error)
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

func (s *ComponentService) GetAllComponents(c *gin.Context) ([]model.ComponentResponse, int, *common.Error) {
	var query model.ComponentFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	components, total, err := s.componentRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	res := make([]model.ComponentResponse, len(components))
	for i, comp := range components {
		resp := modelToComponentResponse(&comp)
		resp.ComponentCategory = s.getComponentCategories(int(comp.ComponentID))
		// Quantity can be joined from inventory, but for simplicty we set 0 if not handled by repo
		res[i] = resp
	}

	return res, total, nil
}

func (s *ComponentService) GetComponentById(c *gin.Context) (*model.ComponentResponse, *common.Error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	comp, err := s.componentRepo.GetByComponentId(id)
	if err != nil || comp == nil {
		return nil, common.NotFound
	}

	res := modelToComponentResponse(comp)
	res.ComponentCategory = s.getComponentCategories(id)
	return &res, nil
}

func (s *ComponentService) CreateComponent(c *gin.Context) (*model.ComponentResponse, *common.Error) {
	var req model.ComponentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	comp := &model.Component{
		ComponentName: req.ComponentName,
		Description:   req.Description,
		Unit:          req.Unit,
		UnitPrice:     req.UnitPrice,
		CreatedAt:     time.Now(),
	}

	categories := make([]model.ComponentCategoryMap, 0)
	for _, cat := range req.ComponentCategory {
		mapEntry := model.ComponentCategoryMap{
			CategoryID: int(cat.CategoryID),
		}
		categories = append(categories, mapEntry)
	}

	result, err := s.componentRepo.CreateComponentTx(comp, categories)
	if err != nil {
		return nil, common.SystemError
	}

	res := modelToComponentResponse(result)
	res.ComponentCategory = s.getComponentCategories(int(result.ComponentID))
	return &res, nil
}

func (s *ComponentService) UpdateComponent(c *gin.Context) *common.Error {
	var req model.ComponentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return &common.Error{Code: "400", Message: err.Error()}
	}

	comp, err := s.componentRepo.GetByComponentId(req.ComponentID)
	if err != nil || comp == nil {
		return common.NotFound
	}

	if req.ComponentName != "" {
		comp.ComponentName = req.ComponentName
	}
	if req.Description != "" {
		comp.Description = req.Description
	}
	if req.Unit != "" {
		comp.Unit = req.Unit
	}
	if req.UnitPrice != 0 {
		comp.UnitPrice = req.UnitPrice
	}
	if req.UpdatedBy != 0 {
		comp.UpdatedBy = int(req.UpdatedBy)
	}
	comp.UpdatedAt = time.Now()

	categories := make([]model.ComponentCategoryMap, 0)
	for _, cat := range req.ComponentCategory {
		mapEntry := model.ComponentCategoryMap{
			CategoryID: int(cat.CategoryID),
		}
		categories = append(categories, mapEntry)
	}

	if err := s.componentRepo.UpdateComponentTx(comp, categories); err != nil {
		return common.SystemError
	}

	return nil
}

func (s *ComponentService) DeleteComponent(c *gin.Context) *common.Error {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	if err := s.componentRepo.DeleteComponentTx(ids); err != nil {
		return common.SystemError
	}

	return nil
}

func (s *ComponentService) getComponentCategories(componentId int) []model.ComponentCategoryDTO {
	var categories []model.ComponentCategoryDTO
	s.db.Model(&model.ComponentCategory{}).
		Select("component_category.category_id, component_category.category_name").
		Joins("INNER JOIN component_category_map ON component_category.category_id = component_category_map.category_id").
		Where("component_category_map.component_id = ?", componentId).
		Scan(&categories)
	return categories
}

func modelToComponentResponse(c *model.Component) model.ComponentResponse {
	return model.ComponentResponse{
		ComponentID:   int(c.ComponentID),
		ComponentName: c.ComponentName,
		Description:   c.Description,
		Unit:          c.Unit,
		UnitPrice:     c.UnitPrice,
		CreatedBy:     int(c.CreatedBy),
		CreatedAt:     c.CreatedAt,
		UpdatedBy:     int(c.UpdatedBy),
		UpdatedAt:     c.UpdatedAt,
	}
}
