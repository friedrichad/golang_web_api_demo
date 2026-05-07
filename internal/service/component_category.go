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
)

type IComponentCategoryService interface {
	GetAllComponentCategories(c *gin.Context) ([]dtos.ComponentCategoryResponse, int, *common.Error)
	GetComponentCategoryById(c *gin.Context) (*dtos.ComponentCategoryResponse, *common.Error)
	CreateComponentCategory(c *gin.Context) (*dtos.ComponentCategoryResponse, *common.Error)
	UpdateComponentCategory(c *gin.Context) *common.Error
	DeleteComponentCategory(c *gin.Context) *common.Error
}

type ComponentCategoryService struct {
	categoryRepo repository.IComponentCategoryRepository
}

var componentCategoryService IComponentCategoryService

func NewComponentCategoryService() IComponentCategoryService {
	if componentCategoryService == nil {
		componentCategoryService = &ComponentCategoryService{
			categoryRepo: repository.NewComponentCategoryRepository(),
		}
	}
	return componentCategoryService
}

func (s *ComponentCategoryService) GetAllComponentCategories(c *gin.Context) ([]dtos.ComponentCategoryResponse, int, *common.Error) {
	var query dtos.ComponentCategoryFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	categories, total, err := s.categoryRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	categoryResponses := make([]dtos.ComponentCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = modelToComponentCategoryResponse(&category)
	}

	return categoryResponses, total, nil
}

func (s *ComponentCategoryService) GetComponentCategoryById(c *gin.Context) (*dtos.ComponentCategoryResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	categoryId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	category, err := s.categoryRepo.GetById(int(categoryId))
	if err != nil {
		return nil, common.NotFound
	}

	if category == nil {
		return nil, &common.Error{Code: "404", Message: "Danh mục thành phần không tồn tại"}
	}

	categoryResponse := modelToComponentCategoryResponse(category)
	return &categoryResponse, nil
}

func (s *ComponentCategoryService) CreateComponentCategory(c *gin.Context) (*dtos.ComponentCategoryResponse, *common.Error) {
	var req dtos.ComponentCategoryCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	categoryRepoTx := s.categoryRepo.(*repository.ComponentCategoryRepository).WithTx(tx)

	category := &model.ComponentCategory{
		CategoryName: req.CategoryName,
		CreatedAt:    time.Now(),
	}

	err := categoryRepoTx.Save(category)
	if err != nil {
		tx.Rollback()
		return nil, common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, common.SystemError
	}

	categoryResponse := modelToComponentCategoryResponse(category)
	return &categoryResponse, nil
}

func (s *ComponentCategoryService) UpdateComponentCategory(c *gin.Context) *common.Error {
	var req dtos.ComponentCategoryUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	categoryRepoTx := s.categoryRepo.(*repository.ComponentCategoryRepository).WithTx(tx)

	category, err := categoryRepoTx.GetById(int(req.CategoryID))
	if err != nil {
		tx.Rollback()
		return common.NotFound
	}

	if category == nil {
		tx.Rollback()
		return &common.Error{Code: "404", Message: "Danh mục thành phần không tồn tại"}
	}

	if req.CategoryName != "" {
		category.CategoryName = req.CategoryName
	}
	category.UpdatedBy = int(req.UpdatedBy)
	category.UpdatedAt = time.Now()

	err = categoryRepoTx.Update(category)
	if err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func (s *ComponentCategoryService) DeleteComponentCategory(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		categoryId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(categoryId)
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	categoryRepoTx := s.categoryRepo.(*repository.ComponentCategoryRepository).WithTx(tx)

	err := categoryRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func modelToComponentCategoryResponse(category *model.ComponentCategory) dtos.ComponentCategoryResponse {
	return dtos.ComponentCategoryResponse{
		CategoryID:   int(category.CategoryID),
		CategoryName: category.CategoryName,
		CreatedBy:    int(category.CreatedBy),
		CreatedAt:    category.CreatedAt,
		UpdatedBy:    int(category.UpdatedBy),
		UpdatedAt:    category.UpdatedAt,
	}
}
