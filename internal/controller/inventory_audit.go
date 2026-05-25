package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
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
	BaseController[model.InventoryAuditResponse]
	auditService service.IInventoryAuditService
}

func NewInventoryAuditController() IInventoryAuditController {
	return &InventoryAuditController{
		auditService: service.NewInventoryAuditService(),
	}
}

// GetAllAudits godoc
// @Summary Get all inventory audits
// @Description Get a list of all inventory audits with optional filtering and pagination
// @Tags audits
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.InventoryAuditResponse
// @Failure 400 {object} common.Error
// @Router /audits [get]
func (c *InventoryAuditController) GetAllAudits() gin.HandlerFunc {
	return c.ResponsePage(c.auditService.GetAllInventoryAudits)
}

// GetAuditById godoc
// @Summary Get an audit by ID
// @Description Retrieve details of a specific inventory audit by its ID
// @Tags audits
// @Accept json
// @Produce json
// @Param id path int true "Audit ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.InventoryAuditResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /audits/{id} [get]
func (c *InventoryAuditController) GetAuditById() gin.HandlerFunc {
	return c.ResponsePointer(c.auditService.GetInventoryAuditById)
}

// CreateAudit godoc
// @Summary Create a new inventory audit
// @Description Add a new inventory audit to the system
// @Tags audits
// @Accept json
// @Produce json
// @Param request body model.InventoryAuditCreate true "Audit Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.InventoryAuditResponse
// @Failure 400 {object} common.Error
// @Router /audits [post]
func (c *InventoryAuditController) CreateAudit() gin.HandlerFunc {
	return c.ResponsePointer(c.auditService.CreateInventoryAudit)
}

// UpdateAudit godoc
// @Summary Update an existing inventory audit
// @Description Modify details of an existing inventory audit
// @Tags audits
// @Accept json
// @Produce json
// @Param request body model.InventoryAuditUpdate true "Audit Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /audits [put]
func (c *InventoryAuditController) UpdateAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.UpdateInventoryAudit)
}

func (c *InventoryAuditController) DeleteAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.DeleteInventoryAudit)
}

// ApproveAudit godoc
// @Summary Approve an inventory audit
// @Description Approve a specific inventory audit
// @Tags audits
// @Accept json
// @Produce json
// @Param request body model.ApprovalAudit true "Approval Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /audits/approval [post]
func (c *InventoryAuditController) ApproveAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.ApprovalAudit)
}

// ConfirmAudit godoc
// @Summary Confirm an inventory audit
// @Description Confirm a specific inventory audit
// @Tags audits
// @Accept json
// @Produce json
// @Param request body model.ConfirmAudit true "Confirmation Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /audits/confirm [post]
func (c *InventoryAuditController) ConfirmAudit() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.auditService.ConfirmAudit)
}
