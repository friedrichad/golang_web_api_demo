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

// GetAllComponents godoc
// @Summary Get all components
// @Description Get a list of all components with optional filtering and pagination
// @Tags components
// @Accept json
// @Produce json
// @Param component_name query string false "Component Name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.ComponentResponse
// @Failure 400 {object} common.Error
// @Router /components [get]
func (c *ComponentController) GetAllComponents() gin.HandlerFunc {
	return c.ResponsePage(c.componentService.GetAllComponents)
}

// GetComponentById godoc
// @Summary Get a component by ID
// @Description Retrieve details of a specific component by its ID
// @Tags components
// @Accept json
// @Produce json
// @Param id path int true "Component ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.ComponentResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /components/{id} [get]
func (c *ComponentController) GetComponentById() gin.HandlerFunc {
	return c.ResponsePointer(c.componentService.GetComponentById)
}

// CreateComponent godoc
// @Summary Create a new component
// @Description Add a new component to the system
// @Tags components
// @Accept json
// @Produce json
// @Param request body model.ComponentCreate true "Component Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.ComponentResponse
// @Failure 400 {object} common.Error
// @Router /components [post]
func (c *ComponentController) CreateComponent() gin.HandlerFunc {
	return c.ResponsePointer(c.componentService.CreateComponent)
}

// UpdateComponent godoc
// @Summary Update an existing component
// @Description Modify details of an existing component
// @Tags components
// @Accept json
// @Produce json
// @Param request body model.ComponentUpdate true "Component Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /components [put]
func (c *ComponentController) UpdateComponent() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.componentService.UpdateComponent)
}

// DeleteComponent godoc
// @Summary Delete components
// @Description Delete one or multiple components by their IDs
// @Tags components
// @Accept json
// @Produce json
// @Param request body []int true "Array of Component IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /components [delete]
func (c *ComponentController) DeleteComponent() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.componentService.DeleteComponent)
}
