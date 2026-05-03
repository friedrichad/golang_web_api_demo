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

type IInventoryAdjustmentService interface {
	GetAllInventoryAdjustments(c *gin.Context) ([]dtos.InventoryAdjustmentResponse, int, *common.Error)
	GetInventoryAdjustmentById(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error)
	CreateInventoryAdjustment(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error)
	UpdateInventoryAdjustment(c *gin.Context) *common.Error
	DeleteInventoryAdjustment(c *gin.Context) *common.Error
}

type InventoryAdjustmentService struct {
	adjustmentRepo repository.IInventoryAdjustment
}

var inventoryAdjustmentService IInventoryAdjustmentService

func NewInventoryAdjustmentService() IInventoryAdjustmentService {
	if inventoryAdjustmentService == nil {
		inventoryAdjustmentService = &InventoryAdjustmentService{
			adjustmentRepo: repository.NewInventoryAdjustmentRepository(),
		}
	}
	return inventoryAdjustmentService
}

func (s *InventoryAdjustmentService) GetAllInventoryAdjustments(c *gin.Context) ([]dtos.InventoryAdjustmentResponse, int, *common.Error) {
	var query dtos.InventoryAdjustmentFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	adjustments, total, err := s.adjustmentRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	adjustmentResponses := make([]dtos.InventoryAdjustmentResponse, len(adjustments))
	for i, adjustment := range adjustments {
		adjustmentResponses[i] = modelToInventoryAdjustmentResponse(&adjustment)
	}

	return adjustmentResponses, total, nil
}

func (s *InventoryAdjustmentService) GetInventoryAdjustmentById(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	adjustmentId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return nil, common.RequestInvalid
	}

	adjustment, err := s.adjustmentRepo.GetById(int32(adjustmentId))
	if err != nil {
		return nil, common.NotFound
	}

	if adjustment == nil {
		return nil, &common.Error{Code: "404", Message: "Điều chỉnh không tồn tại"}
	}

	adjustmentResponse := modelToInventoryAdjustmentResponse(adjustment)
	return &adjustmentResponse, nil
}

func (s *InventoryAdjustmentService) CreateInventoryAdjustment(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error) {
	var req dtos.InventoryAdjustmentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	adjustment := &model.InventoryAdjustment{
		AuditID:     int32(req.AuditID),
		Description: req.Description,
		Note:        req.Note,
		StatusInt:   1,
		CreatedAt:   time.Now(),
	}

	err := s.adjustmentRepo.Save(adjustment)
	if err != nil {
		return nil, common.SystemError
	}

	adjustmentResponse := modelToInventoryAdjustmentResponse(adjustment)
	return &adjustmentResponse, nil
}

func (s *InventoryAdjustmentService) UpdateInventoryAdjustment(c *gin.Context) *common.Error {
	var req dtos.InventoryAdjustmentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	adjustment, err := s.adjustmentRepo.GetById(int32(req.AdjustmentID))
	if err != nil {
		return common.NotFound
	}

	if adjustment == nil {
		return &common.Error{Code: "404", Message: "Điều chỉnh không tồn tại"}
	}

	if req.ApprovedID != 0 {
		adjustment.ApprovedID = int32(req.ApprovedID)
	}
	if req.Description != "" {
		adjustment.Description = req.Description
	}
	if !req.ApprovedTime.IsZero() {
		adjustment.ApprovedTime = req.ApprovedTime
	}
	if req.StatusInt != 0 {
		adjustment.StatusInt = int32(req.StatusInt)
	}
	if req.Note != "" {
		adjustment.Note = req.Note
	}
	adjustment.UpdatedBy = int32(req.UpdatedBy)
	adjustment.UpdatedAt = time.Now()

	err = s.adjustmentRepo.Update(adjustment)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAdjustmentService) DeleteInventoryAdjustment(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int32, len(idStrs))
	for i, idStr := range idStrs {
		adjustmentId, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int32(adjustmentId)
	}

	err := s.adjustmentRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func modelToInventoryAdjustmentResponse(adjustment *model.InventoryAdjustment) dtos.InventoryAdjustmentResponse {
	return dtos.InventoryAdjustmentResponse{
		AdjustmentID: int(adjustment.AdjustmentID),
		AuditID:      int(adjustment.AuditID),
		ApprovedID:   int(adjustment.ApprovedID),
		Description:  adjustment.Description,
		ApprovedTime: adjustment.ApprovedTime,
		StatusInt:    int(adjustment.StatusInt),
		Note:         adjustment.Note,
		CreatedBy:    int(adjustment.CreatedBy),
		CreatedAt:    adjustment.CreatedAt,
	}
}
