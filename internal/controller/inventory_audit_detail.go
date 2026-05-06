package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
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
	BaseController[dtos.InventoryAuditDetailResponse]
	auditService service.IInventoryAuditService
}

func NewInventoryAuditDetailController() IInventoryAuditDetailController {
	return &InventoryAuditDetailController{
		auditService: service.NewInventoryAuditService(),
	}
}
func (c *InventoryAuditDetailController) GetAllInventoryAuditDetails() gin.HandlerFunc {
	return c.ResponsePage(c.auditService.GetAllInventoryAuditDetails)
}
func (c *InventoryAuditDetailController) CreateInventoryAuditDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.CreateInventoryAuditDetail)
}
func (c *InventoryAuditDetailController) UpdateInventoryAuditDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.UpdateInventoryAuditDetail)
}
func (c *InventoryAuditDetailController) DeleteInventoryAuditDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.DeleteInventoryAuditDetail)
}
