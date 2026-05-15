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

func (c *ComponentCategoryController) GetAllCategories() gin.HandlerFunc {
	return c.ResponsePage(c.categoryService.GetAllComponentCategories)
}

func (c *ComponentCategoryController) GetCategoryById() gin.HandlerFunc {
	return c.ResponsePointer(c.categoryService.GetComponentCategoryById)
}

func (c *ComponentCategoryController) CreateCategory() gin.HandlerFunc {
	return c.ResponsePointer(c.categoryService.CreateComponentCategory)
}

func (c *ComponentCategoryController) UpdateCategory() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.categoryService.UpdateComponentCategory)
}

func (c *ComponentCategoryController) DeleteCategory() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.categoryService.DeleteComponentCategory)
}
