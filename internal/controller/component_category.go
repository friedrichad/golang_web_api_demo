package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IComponentCategoryController interface {
	GetAllCategories() gin.HandlerFunc
	GetCategoryById() gin.HandlerFunc
	CreateCategory() gin.HandlerFunc
	UpdateCategory() gin.HandlerFunc
	DeleteCategory() gin.HandlerFunc
}

type ComponentCategoryController struct {
	BaseController[model.ComponentCategoryResponse]
	categoryService service.IComponentCategoryService
}

func NewComponentCategoryController() IComponentCategoryController {
	return &ComponentCategoryController{
		categoryService: service.NewComponentCategoryService(),
	}
}

// GetAllCategories godoc
// @Summary Get all component categories
// @Description Get a list of all component categories with optional filtering and pagination
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.ComponentCategoryResponse
// @Failure 400 {object} common.Error
// @Router /categories [get]
func (c *ComponentCategoryController) GetAllCategories() gin.HandlerFunc {
	return c.ResponsePage(c.categoryService.GetAllComponentCategories)
}

// GetCategoryById godoc
// @Summary Get a category by ID
// @Description Retrieve details of a specific component category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.ComponentCategoryResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /categories/{id} [get]
func (c *ComponentCategoryController) GetCategoryById() gin.HandlerFunc {
	return c.ResponsePointer(c.categoryService.GetComponentCategoryById)
}

// CreateCategory godoc
// @Summary Create a new component category
// @Description Add a new component category to the system
// @Tags categories
// @Accept json
// @Produce json
// @Param request body model.ComponentCategoryCreate true "Category Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.ComponentCategoryResponse
// @Failure 400 {object} common.Error
// @Router /categories [post]
func (c *ComponentCategoryController) CreateCategory() gin.HandlerFunc {
	return c.ResponsePointer(c.categoryService.CreateComponentCategory)
}

// UpdateCategory godoc
// @Summary Update an existing component category
// @Description Modify details of an existing component category
// @Tags categories
// @Accept json
// @Produce json
// @Param request body model.ComponentCategoryUpdate true "Category Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /categories [put]
func (c *ComponentCategoryController) UpdateCategory() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.categoryService.UpdateComponentCategory)
}

// DeleteCategory godoc
// @Summary Delete component categories
// @Description Delete one or multiple component categories by their IDs
// @Tags categories
// @Accept json
// @Produce json
// @Param request body []int true "Array of Category IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /categories [delete]
func (c *ComponentCategoryController) DeleteCategory() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.categoryService.DeleteComponentCategory)
}
