package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/middleware"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IInventoryAdjustmentService interface {
	GetAllInventoryAdjustments(c *gin.Context) ([]dtos.InventoryAdjustmentResponse, int, *common.Error)
	GetInventoryAdjustmentById(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error)
	CreateInventoryAdjustment(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error)
	UpdateInventoryAdjustment(c *gin.Context) *common.Error
	DeleteInventoryAdjustment(c *gin.Context) *common.Error
	ApproveInventoryAdjustment(c *gin.Context) *common.Error
	GetAllInventoryAdjustmentDetails(c *gin.Context) ([]dtos.InventoryAdjustmentDetailResponse, int, *common.Error)
	GetInventoryAdjustmentDetailById(c *gin.Context) (*dtos.InventoryAdjustmentDetailResponse, *common.Error)
	CreateInventoryAdjustmentDetail(c *gin.Context) (*dtos.InventoryAdjustmentDetailResponse, *common.Error)
	UpdateInventoryAdjustmentDetail(c *gin.Context) *common.Error
	DeleteInventoryAdjustmentDetail(c *gin.Context) *common.Error
}

type InventoryAdjustmentService struct {
	adjustmentRepo   repository.IInventoryAdjustment
	detailRepo       repository.IInventoryAdjustmentDetail
	componentBinRepo repository.IComponentBin
	ledgerService    IInventoryLedgerService
}

var inventoryAdjustmentService IInventoryAdjustmentService

func NewInventoryAdjustmentService() IInventoryAdjustmentService {
	if inventoryAdjustmentService == nil {
		inventoryAdjustmentService = &InventoryAdjustmentService{
			adjustmentRepo:   repository.NewInventoryAdjustmentRepository(),
			detailRepo:       repository.NewInventoryAdjustmentDetailRepository(),
			componentBinRepo: repository.NewComponentBinRepository(),
			ledgerService:    NewInventoryLedgerService(),
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

	adjustmentId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	adjustment, err := s.adjustmentRepo.GetByAdjustmentId(int(adjustmentId))
	if err != nil {
		return nil, common.NotFound
	}

	if adjustment == nil {
		return nil, &common.Error{Code: "404", Message: "Điều chỉnh không tồn tại"}
	}

	adjustmentResponse := modelToInventoryAdjustmentResponse(adjustment)

	details, err := s.detailRepo.GetByAdjustmentId(adjustment.AdjustmentID)
	if err == nil && len(details) > 0 {
		adjustmentResponse.Details = make([]dtos.InventoryAdjustmentDetailResponse, len(details))
		for i, detail := range details {
			adjustmentResponse.Details[i] = modelToInventoryAdjustmentDetailResponse(&detail)
		}
	}

	return &adjustmentResponse, nil
}

func (s *InventoryAdjustmentService) CreateInventoryAdjustment(c *gin.Context) (*dtos.InventoryAdjustmentResponse, *common.Error) {
	var req dtos.InventoryAdjustmentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	adjustment := &model.InventoryAdjustment{
		AuditID:     int(req.AuditID),
		Description: req.Description,
		Note:        req.Note,
		StatusInt:   constants.InventoryAdjustmentStatusPending,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
	}

	details := make([]model.InventoryAdjustmentDetail, 0)
	for _, detailReq := range req.Details {
		detail := model.InventoryAdjustmentDetail{
			ComponentID:        detailReq.ComponentID,
			BinID:              detailReq.BinID,
			WarehouseID:        detailReq.WarehouseID,
			QuantityBefore:     detailReq.QuantityBefore,
			QuantityAfter:      detailReq.QuantityAfter,
			AdjustmentQuantity: detailReq.AdjustmentQuantity,
			CreatedBy:          userID,
			CreatedAt:          time.Now(),
		}
		details = append(details, detail)
	}

	result, err := s.adjustmentRepo.CreateInventoryAdjustmentTx(adjustment, details)
	if err != nil {
		return nil, common.SystemError
	}

	adjustmentResponse := modelToInventoryAdjustmentResponse(result)
	return &adjustmentResponse, nil
}

func (s *InventoryAdjustmentService) UpdateInventoryAdjustment(c *gin.Context) *common.Error {
	var req dtos.InventoryAdjustmentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	adjustment, err := s.adjustmentRepo.GetByAdjustmentId(int(req.AdjustmentID))
	if err != nil || adjustment == nil {
		return common.NotFound
	}

	if adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
		return &common.Error{Code: "403", Message: "Chỉ được phép chỉnh sửa đơn ở trạng thái chờ duyệt"}
	}

	if req.Description != "" {
		adjustment.Description = req.Description
	}
	if req.Note != "" {
		adjustment.Note = req.Note
	}
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	adjustment.UpdatedBy = userID
	adjustment.UpdatedAt = time.Now()

	if err := s.adjustmentRepo.UpdateInventoryAdjustmentTx(adjustment); err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAdjustmentService) DeleteInventoryAdjustment(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		adjustmentId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(adjustmentId)
	}

	// Validate all items have pending status
	for _, id := range ids {
		adjustment, err := s.adjustmentRepo.GetByAdjustmentId(id)
		if err != nil || adjustment == nil {
			continue
		}
		if adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
			return &common.Error{Code: "403", Message: "Chỉ được xóa đơn ở trạng thái chờ duyệt"}
		}
	}

	if err := s.adjustmentRepo.DeleteInventoryAdjustmentTx(ids); err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAdjustmentService) ApproveInventoryAdjustment(c *gin.Context) *common.Error {
	var req struct {
		AdjustmentID int `json:"adjustment_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	adjustment, err := s.adjustmentRepo.GetByAdjustmentId(req.AdjustmentID)
	if err != nil || adjustment == nil {
		return common.NotFound
	}

	if adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
		return &common.Error{Code: "403", Message: "Đơn đã được xử lý"}
	}

	// Update adjustment status
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	adjustment.StatusInt = constants.InventoryAdjustmentStatusApproved
	adjustment.ApprovedID = userID
	adjustment.UpdatedAt = time.Now()
	adjustment.UpdatedBy = userID

	// Get details
	details, err := s.detailRepo.GetByAdjustmentId(adjustment.AdjustmentID)
	if err != nil {
		return common.SystemError
	}

	// Update component bins and create ledger entries
	binRepo := repository.NewBinRepository()

	for _, detail := range details {
		compBin, err := s.componentBinRepo.GetByComponentAndBinId(detail.ComponentID, detail.BinID)
		if err != nil {
			return common.SystemError
		}

		if compBin == nil {
			compBin = &model.ComponentBin{
				ComponentID: detail.ComponentID,
				BinID:       detail.BinID,
				Quantity:    detail.QuantityAfter,
				CreatedAt:   time.Now(),
				CreatedBy:   userID,
			}
			err = s.componentBinRepo.Save(compBin)
		} else {
			compBin.Quantity = detail.QuantityAfter
			compBin.UpdatedAt = time.Now()
			compBin.UpdatedBy = userID
			err = s.componentBinRepo.Update(compBin)
		}

		if err != nil {
			return common.SystemError
		}

		// Create ledger entry
		binInfo, _ := binRepo.GetById(detail.BinID)
		warehouseID := 0
		if binInfo != nil {
			warehouseID = binInfo.WarehouseID
		}

		quantityChange := detail.QuantityAfter - detail.QuantityBefore
		ledgerReq := &dtos.InventoryLedgerCreate{
			ComponentID:     detail.ComponentID,
			WarehouseID:     warehouseID,
			BinID:           detail.BinID,
			ReferenceType:   constants.LedgerReferenceTypeAdjustment,
			ReferenceTypeID: adjustment.AdjustmentID,
			Description:     fmt.Sprintf("Điều chỉnh từ khoá #%d", adjustment.AdjustmentID),
			QuantityChange:  quantityChange,
			QuantityAfter:   detail.QuantityAfter,
			Note:            adjustment.Note,
			CreatedBy:       userID,
		}
		if err := s.ledgerService.CreateInventoryLedgerEntry(ledgerReq); err != nil {
			return &common.Error{Code: "500", Message: fmt.Sprintf("Lỗi tạo sổ cái: %v", err)}
		}
	}

	// Update adjustment status in repository
	if err := s.adjustmentRepo.Update(adjustment); err != nil {
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
		StatusInt:    int(adjustment.StatusInt),
		Note:         adjustment.Note,
		CreatedBy:    int(adjustment.CreatedBy),
		CreatedAt:    adjustment.CreatedAt,
		UpdatedBy:    int(adjustment.UpdatedBy),
		UpdatedAt:    adjustment.UpdatedAt,
	}
}

func modelToInventoryAdjustmentDetailResponse(detail *model.InventoryAdjustmentDetail) dtos.InventoryAdjustmentDetailResponse {
	return dtos.InventoryAdjustmentDetailResponse{
		AdjustmentDetailID: int(detail.AdjustmentDetailID),
		AdjustmentID:       int(detail.AdjustmentID),
		ComponentID:        int(detail.ComponentID),
		BinID:              int(detail.BinID),
		WarehouseID:        int(detail.WarehouseID),
		QuantityBefore:     detail.QuantityBefore,
		QuantityAfter:      detail.QuantityAfter,
		AdjustmentQuantity: detail.AdjustmentQuantity,
		CreatedBy:          int(detail.CreatedBy),
		CreatedAt:          detail.CreatedAt,
		UpdatedBy:          int(detail.UpdatedBy),
		UpdatedAt:          detail.UpdatedAt,
	}
}

// InventoryAdjustmentDetail CRUD operations

func (s *InventoryAdjustmentService) GetAllInventoryAdjustmentDetails(c *gin.Context) ([]dtos.InventoryAdjustmentDetailResponse, int, *common.Error) {
	var query dtos.InventoryAdjustmentDetailFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	details, total, err := s.detailRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	detailResponses := make([]dtos.InventoryAdjustmentDetailResponse, len(details))
	for i, detail := range details {
		detailResponses[i] = modelToInventoryAdjustmentDetailResponse(&detail)
	}

	return detailResponses, total, nil
}

func (s *InventoryAdjustmentService) GetInventoryAdjustmentDetailById(c *gin.Context) (*dtos.InventoryAdjustmentDetailResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	detailId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	detail, err := s.detailRepo.GetByAdjustmentDetailId(int(detailId))
	if err != nil {
		return nil, common.NotFound
	}

	if detail == nil {
		return nil, &common.Error{Code: "404", Message: "Chi tiết điều chỉnh không tồn tại"}
	}

	detailResponse := modelToInventoryAdjustmentDetailResponse(detail)
	return &detailResponse, nil
}

func (s *InventoryAdjustmentService) CreateInventoryAdjustmentDetail(c *gin.Context) (*dtos.InventoryAdjustmentDetailResponse, *common.Error) {
	var req dtos.InventoryAdjustmentDetailCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	adjustment, err := s.adjustmentRepo.GetByAdjustmentId(req.AdjustmentID)
	if err != nil || adjustment == nil {
		return nil, &common.Error{Code: "404", Message: "Điều chỉnh không tồn tại"}
	}

	if adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
		return nil, &common.Error{Code: "403", Message: "Chỉ được phép thêm chi tiết cho đơn ở trạng thái chờ duyệt"}
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))

	detail := &model.InventoryAdjustmentDetail{
		AdjustmentID:       req.AdjustmentID,
		ComponentID:        req.ComponentID,
		BinID:              req.BinID,
		WarehouseID:        req.WarehouseID,
		QuantityBefore:     req.QuantityBefore,
		QuantityAfter:      req.QuantityAfter,
		AdjustmentQuantity: req.AdjustmentQuantity,
		CreatedBy:          userID,
		CreatedAt:          time.Now(),
	}

	err = s.detailRepo.Save(detail)
	if err != nil {
		return nil, common.SystemError
	}

	detailResponse := modelToInventoryAdjustmentDetailResponse(detail)
	return &detailResponse, nil
}

func (s *InventoryAdjustmentService) UpdateInventoryAdjustmentDetail(c *gin.Context) *common.Error {
	var req dtos.InventoryAdjustmentDetailUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	detail, err := s.detailRepo.GetByAdjustmentDetailId(req.AdjustmentDetailID)
	if err != nil || detail == nil {
		return &common.Error{Code: "404", Message: "Chi tiết điều chỉnh không tồn tại"}
	}

	adjustment, err := s.adjustmentRepo.GetByAdjustmentId(detail.AdjustmentID)
	if err != nil || adjustment == nil {
		return &common.Error{Code: "404", Message: "Điều chỉnh không tồn tại"}
	}

	if adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
		return &common.Error{Code: "403", Message: "Chỉ được phép chỉnh sửa chi tiết của đơn ở trạng thái chờ duyệt"}
	}

	if req.ComponentID != 0 {
		detail.ComponentID = req.ComponentID
	}
	if req.BinID != 0 {
		detail.BinID = req.BinID
	}
	if req.WarehouseID != 0 {
		detail.WarehouseID = req.WarehouseID
	}
	detail.QuantityBefore = req.QuantityBefore
	detail.QuantityAfter = req.QuantityAfter
	detail.AdjustmentQuantity = req.AdjustmentQuantity
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	detail.UpdatedBy = userID
	detail.UpdatedAt = time.Now()

	err = s.detailRepo.Update(detail)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAdjustmentService) DeleteInventoryAdjustmentDetail(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		detailId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(detailId)
	}

	// Check if adjustment is pending for all details
	for _, id := range ids {
		detail, err := s.detailRepo.GetByAdjustmentDetailId(id)
		if err != nil || detail == nil {
			continue
		}
		adjustment, err := s.adjustmentRepo.GetByAdjustmentId(detail.AdjustmentID)
		if err == nil && adjustment != nil {
			if adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
				return &common.Error{Code: "403", Message: "Chỉ được xóa chi tiết của đơn ở trạng thái chờ duyệt"}
			}
		}
	}

	err := s.detailRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}
