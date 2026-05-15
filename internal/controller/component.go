package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IComponentController interface {
	GetAllComponents() gin.HandlerFunc
	GetComponentById() gin.HandlerFunc
	CreateComponent() gin.HandlerFunc
	UpdateComponent() gin.HandlerFunc
	DeleteComponent() gin.HandlerFunc
}

type ComponentController struct {
	BaseController[model.ComponentResponse]
	componentService service.IComponentService
}

func NewComponentController() IComponentController {
	componentService := service.NewComponentService()
	return &ComponentController{componentService: componentService}
}

func (c *ComponentController) GetAllComponents() gin.HandlerFunc {
	return c.ResponsePage(c.componentService.GetAllComponents)
}

func (c *ComponentController) GetComponentById() gin.HandlerFunc {
	return c.ResponsePointer(c.componentService.GetComponentById)
}

func (c *ComponentController) CreateComponent() gin.HandlerFunc {
	return c.ResponsePointer(c.componentService.CreateComponent)
}

func (c *ComponentController) UpdateComponent() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.componentService.UpdateComponent)
}

func (c *ComponentController) DeleteComponent() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.componentService.DeleteComponent)
}
