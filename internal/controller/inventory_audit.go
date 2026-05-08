package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IInventoryAuditController interface {
	GetAllAudits() gin.HandlerFunc
	GetAuditById() gin.HandlerFunc
	CreateAudit() gin.HandlerFunc
	UpdateAudit() gin.HandlerFunc
	DeleteAudit() gin.HandlerFunc
	ApproveAudit() gin.HandlerFunc
	ConfirmAudit() gin.HandlerFunc
}

type InventoryAuditController struct {
	BaseController[dtos.InventoryAuditResponse]
	auditService service.IInventoryAuditService
}

func NewInventoryAuditController() IInventoryAuditController {
	return &InventoryAuditController{
		auditService: service.NewInventoryAuditService(),
	}
}

func (c *InventoryAuditController) GetAllAudits() gin.HandlerFunc {
	return c.ResponsePage(c.auditService.GetAllInventoryAudits)
}

func (c *InventoryAuditController) GetAuditById() gin.HandlerFunc {
	return c.ResponsePointer(c.auditService.GetInventoryAuditById)
}

func (c *InventoryAuditController) CreateAudit() gin.HandlerFunc {
	return c.ResponsePointer(c.auditService.CreateInventoryAudit)
}

func (c *InventoryAuditController) UpdateAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.UpdateInventoryAudit)
}

func (c *InventoryAuditController) DeleteAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.DeleteInventoryAudit)
}
func (c *InventoryAuditController) ApproveAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.ApprovalAudit)
}

func (c *InventoryAuditController) ConfirmAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.ConfirmAudit)
}
