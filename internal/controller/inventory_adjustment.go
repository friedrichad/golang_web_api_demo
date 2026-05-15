package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IInventoryAdjustmentController interface {
	GetAllAdjustments() gin.HandlerFunc
	GetAdjustmentById() gin.HandlerFunc
	CreateAdjustment() gin.HandlerFunc
	UpdateAdjustment() gin.HandlerFunc
	DeleteAdjustment() gin.HandlerFunc
	ApproveAdjustment() gin.HandlerFunc
}

type InventoryAdjustmentController struct {
	BaseController[model.InventoryAdjustmentResponse]
	adjustmentService service.IInventoryAdjustmentService
}

func NewInventoryAdjustmentController() IInventoryAdjustmentController {
	return &InventoryAdjustmentController{
		adjustmentService: service.NewInventoryAdjustmentService(),
	}
}

func (c *InventoryAdjustmentController) GetAllAdjustments() gin.HandlerFunc {
	return c.ResponsePage(c.adjustmentService.GetAllInventoryAdjustments)
}

func (c *InventoryAdjustmentController) GetAdjustmentById() gin.HandlerFunc {
	return c.ResponsePointer(c.adjustmentService.GetInventoryAdjustmentById)
}

func (c *InventoryAdjustmentController) CreateAdjustment() gin.HandlerFunc {
	return c.ResponsePointer(c.adjustmentService.CreateInventoryAdjustment)
}

func (c *InventoryAdjustmentController) UpdateAdjustment() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.adjustmentService.UpdateInventoryAdjustment)
}

func (c *InventoryAdjustmentController) DeleteAdjustment() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.adjustmentService.DeleteInventoryAdjustment)
}

func (c *InventoryAdjustmentController) ApproveAdjustment() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.adjustmentService.ApproveInventoryAdjustment)
}
