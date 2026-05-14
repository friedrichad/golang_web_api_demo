package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
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
	BaseController[dtos.PositionResponse]
	positionService service.IPositionService
}

func NewPositionController() IPositionController {
	return &PositionController{
		positionService: service.NewPositionService(),
	}
}

func (c *PositionController) GetAllPositions() gin.HandlerFunc {
	return c.ResponsePage(c.positionService.GetAllPositions)
}

func (c *PositionController) GetPositionById() gin.HandlerFunc {
	return c.ResponsePointer(c.positionService.GetPositionById)
}

func (c *PositionController) CreatePosition() gin.HandlerFunc {
	return c.ResponsePointer(c.positionService.CreatePosition)
}

func (c *PositionController) UpdatePosition() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.positionService.UpdatePosition)
}

func (c *PositionController) DeletePosition() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.positionService.DeletePosition)
}
