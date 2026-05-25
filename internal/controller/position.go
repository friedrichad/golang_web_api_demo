package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IPositionController interface {
	GetAllPositions() gin.HandlerFunc
	GetPositionById() gin.HandlerFunc
	CreatePosition() gin.HandlerFunc
	UpdatePosition() gin.HandlerFunc
	DeletePosition() gin.HandlerFunc
}

type PositionController struct {
	BaseController[model.PositionResponse]
	positionService service.IPositionService
}

func NewPositionController() IPositionController {
	return &PositionController{
		positionService: service.NewPositionService(),
	}
}

// GetAllPositions godoc
// @Summary Get all positions
// @Description Get a list of all positions with optional filtering and pagination
// @Tags positions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.PositionResponse
// @Failure 400 {object} common.Error
// @Router /positions [get]
func (c *PositionController) GetAllPositions() gin.HandlerFunc {
	return c.ResponsePage(c.positionService.GetAllPositions)
}

// GetPositionById godoc
// @Summary Get a position by ID
// @Description Retrieve details of a specific position by its ID
// @Tags positions
// @Accept json
// @Produce json
// @Param id path int true "Position ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.PositionResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /positions/{id} [get]
func (c *PositionController) GetPositionById() gin.HandlerFunc {
	return c.ResponsePointer(c.positionService.GetPositionById)
}

// CreatePosition godoc
// @Summary Create a new position
// @Description Add a new position to the system
// @Tags positions
// @Accept json
// @Produce json
// @Param request body model.PositionCreate true "Position Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.PositionResponse
// @Failure 400 {object} common.Error
// @Router /positions [post]
func (c *PositionController) CreatePosition() gin.HandlerFunc {
	return c.ResponsePointer(c.positionService.CreatePosition)
}

// UpdatePosition godoc
// @Summary Update an existing position
// @Description Modify details of an existing position
// @Tags positions
// @Accept json
// @Produce json
// @Param request body model.PositionUpdate true "Position Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /positions [put]
func (c *PositionController) UpdatePosition() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.positionService.UpdatePosition)
}

// DeletePosition godoc
// @Summary Delete positions
// @Description Delete one or multiple positions by their IDs
// @Tags positions
// @Accept json
// @Produce json
// @Param request body []int true "Array of Position IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /positions [delete]
func (c *PositionController) DeletePosition() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.positionService.DeletePosition)
}
