package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IInventoryAuditDetailController interface {
	GetAllInventoryAuditDetails() gin.HandlerFunc
	CreateInventoryAuditDetail() gin.HandlerFunc
	UpdateInventoryAuditDetail() gin.HandlerFunc
	DeleteInventoryAuditDetail() gin.HandlerFunc
}

type InventoryAuditDetailController struct {
	BaseController[model.InventoryAuditDetailResponse]
	auditService service.IInventoryAuditService
}

func NewInventoryAuditDetailController() IInventoryAuditDetailController {
	return &InventoryAuditDetailController{
		auditService: service.NewInventoryAuditService(),
	}
}

// GetAllInventoryAuditDetails godoc
// @Summary Get all audit details
// @Description Get a list of all inventory audit details with optional filtering and pagination
// @Tags audit-details
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.InventoryAuditDetailResponse
// @Failure 400 {object} common.Error
// @Router /audit-details [get]
func (c *InventoryAuditDetailController) GetAllInventoryAuditDetails() gin.HandlerFunc {
	return c.ResponsePage(c.auditService.GetAllInventoryAuditDetails)
}

// CreateInventoryAuditDetail godoc
// @Summary Create a new audit detail
// @Description Add a new inventory audit detail to an audit
// @Tags audit-details
// @Accept json
// @Produce json
// @Param request body model.InventoryAuditDetailCreate true "Audit Detail Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /audit-details [post]
func (c *InventoryAuditDetailController) CreateInventoryAuditDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.CreateInventoryAuditDetail)
}

// UpdateInventoryAuditDetail godoc
// @Summary Update an existing audit detail
// @Description Modify details of an existing inventory audit detail
// @Tags audit-details
// @Accept json
// @Produce json
// @Param request body model.InventoryAuditDetailUpdate true "Audit Detail Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /audit-details [put]
func (c *InventoryAuditDetailController) UpdateInventoryAuditDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.UpdateInventoryAuditDetail)
}

// DeleteInventoryAuditDetail godoc
// @Summary Delete audit details
// @Description Delete one or multiple inventory audit details by their IDs
// @Tags audit-details
// @Accept json
// @Produce json
// @Param request body []int true "Array of Audit Detail IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /audit-details [delete]
func (c *InventoryAuditDetailController) DeleteInventoryAuditDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.DeleteInventoryAuditDetail)
}
