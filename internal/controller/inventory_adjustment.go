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

// GetAllAdjustments godoc
// @Summary Get all inventory adjustments
// @Description Get a list of all inventory adjustments with optional filtering and pagination
// @Tags adjustments
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.InventoryAdjustmentResponse
// @Failure 400 {object} common.Error
// @Router /adjustments [get]
func (c *InventoryAdjustmentController) GetAllAdjustments() gin.HandlerFunc {
	return c.ResponsePage(c.adjustmentService.GetAllInventoryAdjustments)
}

// GetAdjustmentById godoc
// @Summary Get an adjustment by ID
// @Description Retrieve details of a specific inventory adjustment by its ID
// @Tags adjustments
// @Accept json
// @Produce json
// @Param id path int true "Adjustment ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.InventoryAdjustmentResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /adjustments/{id} [get]
func (c *InventoryAdjustmentController) GetAdjustmentById() gin.HandlerFunc {
	return c.ResponsePointer(c.adjustmentService.GetInventoryAdjustmentById)
}

// CreateAdjustment godoc
// @Summary Create a new inventory adjustment
// @Description Add a new inventory adjustment to the system
// @Tags adjustments
// @Accept json
// @Produce json
// @Param request body model.InventoryAdjustmentCreate true "Adjustment Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.InventoryAdjustmentResponse
// @Failure 400 {object} common.Error
// @Router /adjustments [post]
func (c *InventoryAdjustmentController) CreateAdjustment() gin.HandlerFunc {
	return c.ResponsePointer(c.adjustmentService.CreateInventoryAdjustment)
}

// UpdateAdjustment godoc
// @Summary Update an existing inventory adjustment
// @Description Modify details of an existing inventory adjustment
// @Tags adjustments
// @Accept json
// @Produce json
// @Param request body model.InventoryAdjustmentUpdate true "Adjustment Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /adjustments [put]
func (c *InventoryAdjustmentController) UpdateAdjustment() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.adjustmentService.UpdateInventoryAdjustment)
}

// DeleteAdjustment godoc
// @Summary Delete inventory adjustments
// @Description Delete one or multiple inventory adjustments by their IDs
// @Tags adjustments
// @Accept json
// @Produce json
// @Param request body []int true "Array of Adjustment IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /adjustments [delete]
func (c *InventoryAdjustmentController) DeleteAdjustment() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.adjustmentService.DeleteInventoryAdjustment)
}

// ApproveAdjustment godoc
// @Summary Approve an inventory adjustment
// @Description Approve a specific inventory adjustment
// @Tags adjustments
// @Accept json
// @Produce json
// @Param request body model.InventoryAdjustmentUpdate true "Approval Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /adjustments/approval [post]
func (c *InventoryAdjustmentController) ApproveAdjustment() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.adjustmentService.ApproveInventoryAdjustment)
}
