package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IInventoryAuditService interface {
	GetAllInventoryAudits(c *gin.Context) ([]dtos.InventoryAuditResponse, int, *common.Error)
	GetInventoryAuditById(c *gin.Context) (*dtos.InventoryAuditResponse, *common.Error)
	CreateInventoryAudit(c *gin.Context) (*dtos.InventoryAuditResponse, *common.Error)
	UpdateInventoryAudit(c *gin.Context) *common.Error
	DeleteInventoryAudit(c *gin.Context) *common.Error
}

type InventoryAuditService struct {
	auditRepo repository.IInventoryAudit
}

var inventoryAuditService IInventoryAuditService

func NewInventoryAuditService() IInventoryAuditService {
	if inventoryAuditService == nil {
		inventoryAuditService = &InventoryAuditService{
			auditRepo: repository.NewInventoryAuditRepository(),
		}
	}
	return inventoryAuditService
}

func (s *InventoryAuditService) GetAllInventoryAudits(c *gin.Context) ([]dtos.InventoryAuditResponse, int, *common.Error) {
	var query dtos.InventoryAuditFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	audits, total, err := s.auditRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	auditResponses := make([]dtos.InventoryAuditResponse, len(audits))
	for i, audit := range audits {
		auditResponses[i] = modelToInventoryAuditResponse(&audit)
	}

	return auditResponses, total, nil
}

func (s *InventoryAuditService) GetInventoryAuditById(c *gin.Context) (*dtos.InventoryAuditResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	auditId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return nil, common.RequestInvalid
	}

	audit, err := s.auditRepo.GetById(int32(auditId))
	if err != nil {
		return nil, common.NotFound
	}

	if audit == nil {
		return nil, &common.Error{Code: "404", Message: "Kiểm kê không tồn tại"}
	}

	auditResponse := modelToInventoryAuditResponse(audit)
	return &auditResponse, nil
}

func (s *InventoryAuditService) CreateInventoryAudit(c *gin.Context) (*dtos.InventoryAuditResponse, *common.Error) {
	var req dtos.InventoryAuditCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	audit := &model.InventoryAudit{
		WarehouseID: int32(req.WarehouseID),
		Note:        req.Note,
		StatusInt:   1,
		CreatedAt:   time.Now(),
	}

	err := s.auditRepo.Save(audit)
	if err != nil {
		return nil, common.SystemError
	}

	auditResponse := modelToInventoryAuditResponse(audit)
	return &auditResponse, nil
}

func (s *InventoryAuditService) UpdateInventoryAudit(c *gin.Context) *common.Error {
	var req dtos.InventoryAuditUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	audit, err := s.auditRepo.GetById(int32(req.AuditID))
	if err != nil {
		return common.NotFound
	}

	if audit == nil {
		return &common.Error{Code: "404", Message: "Kiểm kê không tồn tại"}
	}

	if req.StatusInt != 0 {
		audit.StatusInt = int32(req.StatusInt)
	}
	if req.Note != "" {
		audit.Note = req.Note
	}
	audit.UpdatedBy = int32(req.UpdatedBy)
	audit.UpdatedAt = time.Now()

	err = s.auditRepo.Update(audit)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAuditService) DeleteInventoryAudit(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int32, len(idStrs))
	for i, idStr := range idStrs {
		auditId, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int32(auditId)
	}

	err := s.auditRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func modelToInventoryAuditResponse(audit *model.InventoryAudit) dtos.InventoryAuditResponse {
	return dtos.InventoryAuditResponse{
		AuditID:     int(audit.AuditID),
		WarehouseID: int(audit.WarehouseID),
		StatusInt:   int(audit.StatusInt),
		Note:        audit.Note,
		CreatedBy:   int(audit.CreatedBy),
		CreatedAt:   audit.CreatedAt,
		UpdatedBy:   int(audit.UpdatedBy),
		UpdatedAt:   audit.UpdatedAt,
	}
}
