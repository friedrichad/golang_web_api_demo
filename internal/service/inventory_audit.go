package service

import (
	"fmt"
	"strconv"
	"time"

	"math"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/middleware"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IInventoryAuditService interface {
	GetAllInventoryAudits(c *gin.Context) ([]dtos.InventoryAuditResponse, int, *common.Error)
	GetInventoryAuditById(c *gin.Context) (*dtos.InventoryAuditResponse, *common.Error)
	CreateInventoryAudit(c *gin.Context) (*dtos.InventoryAuditResponse, *common.Error)
	UpdateInventoryAudit(c *gin.Context) *common.Error
	DeleteInventoryAudit(c *gin.Context) *common.Error
	ApprovalAudit(c *gin.Context) *common.Error
	ConfirmAudit(c *gin.Context) *common.Error

	CreateInventoryAuditDetail(c *gin.Context) *common.Error
	GetAllInventoryAuditDetails(c *gin.Context) ([]dtos.InventoryAuditDetailResponse, int, *common.Error)
	UpdateInventoryAuditDetail(c *gin.Context) *common.Error
	DeleteInventoryAuditDetail(c *gin.Context) *common.Error
}

type InventoryAuditService struct {
	auditRepo       repository.IInventoryAudit
	auditDetailRepo repository.IInventoryAuditDetail
	ledgerService   IInventoryLedgerService
}

var inventoryAuditService IInventoryAuditService

func NewInventoryAuditService() IInventoryAuditService {
	if inventoryAuditService == nil {
		inventoryAuditService = &InventoryAuditService{
			auditRepo:       repository.NewInventoryAuditRepository(),
			auditDetailRepo: repository.NewInventoryAuditDetailRepository(),
			ledgerService:   NewInventoryLedgerService(),
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

	auditId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	audit, err := s.auditRepo.GetByRequestId(int(auditId))
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

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	audit := &model.InventoryAudit{
		WarehouseID: int(req.WarehouseID),
		Note:        req.Note,
		StatusInt:   1,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
	}

	details := make([]model.InventoryAuditDetail, 0)
	for _, detail := range req.AuditDetail {
		auditDetail := model.InventoryAuditDetail{
			ComponentID:        detail.ComponentID,
			BinID:              detail.BinID,
			SystemQuantity:     detail.SystemQuantity,
			ActualQuantity:     detail.ActualQuantity,
			DifferenceQuantity: math.Abs(detail.SystemQuantity - detail.ActualQuantity),
			CreatedBy:          userID,
		}
		details = append(details, auditDetail)
	}

	result, err := s.auditRepo.CreateInventoryAuditTx(audit, details)
	if err != nil {
		return nil, common.SystemError
	}

	auditResponse := modelToInventoryAuditResponse(result)
	return &auditResponse, nil
}

func (s *InventoryAuditService) UpdateInventoryAudit(c *gin.Context) *common.Error {
	var req dtos.InventoryAuditUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return &common.Error{Code: "400", Message: err.Error()}
	}

	audit, err := s.auditRepo.GetByRequestId(int(req.AuditID))
	if err != nil {
		return common.NotFound
	}

	if audit == nil {
		return &common.Error{Code: "404", Message: "Kiểm kê không tồn tại"}
	}

	if req.StatusInt != 0 {
		audit.StatusInt = int(req.StatusInt)
	}
	if req.Note != "" {
		audit.Note = req.Note
	}
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	audit.UpdatedBy = userID
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

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		auditId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(auditId)
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

func (s *InventoryAuditService) ApprovalAudit(c *gin.Context) *common.Error {
	var req dtos.ApprovalAudit
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	audit, err := s.auditRepo.GetByRequestId(int(req.AuditID))
	if err != nil {
		return common.SystemError
	}

	if audit == nil {
		return common.NotFound
	}

	if !constants.IsValidInventoryAuditStatus(req.StatusInt) {
		return &common.Error{Code: "400", Message: "Trạng thái kiểm kê không hợp lệ!"}
	}

	if audit.StatusInt == constants.InventoryAuditStatusApproved || audit.StatusInt == constants.InventoryAuditStatusRejected {
		return &common.Error{Code: "400", Message: "Không thể phê duyệt lại kiểm kê đã được phê duyệt!"}
	}

	if req.StatusInt != constants.InventoryAuditStatusApproved && req.StatusInt != constants.InventoryAuditStatusRejected {
		return &common.Error{Code: "400", Message: "Kiểm kê chỉ có thể được phê duyệt hoặc từ chối!"}
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	audit.StatusInt = int(req.StatusInt)
	audit.Note = req.Note
	audit.UpdatedBy = userID
	audit.UpdatedAt = time.Now()

	err = s.auditRepo.Update(audit)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAuditService) ConfirmAudit(c *gin.Context) *common.Error {
	var req dtos.ConfirmAudit
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if req.AuditID == 0 {
		return common.RequestInvalid
	}

	audit, err := s.auditRepo.GetByRequestId(int(req.AuditID))
	if err != nil {
		return common.SystemError
	}

	if audit == nil {
		return common.NotFound
	}

	if !constants.IsValidInventoryAuditStatus(req.StatusInt) {
		return &common.Error{Code: "400", Message: "Trạng thái kiểm kê không hợp lệ!"}
	}

	if audit.StatusInt != constants.InventoryAuditStatusApproved {
		return &common.Error{Code: "400", Message: "Chỉ có thể xác nhận kiểm kê đã được phê duyệt!"}
	}

	// If applying the audit (status = APPLIED), update component bins and create ledger entries
	if req.StatusInt == constants.InventoryAuditStatusApplied {
		tx := db.Instance.Begin()
		if tx.Error != nil {
			return common.SystemError
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		compBinRepoTx := repository.NewComponentBinRepository().(*repository.ComponentBinRepository).WithTx(tx)
		auditDetailRepoTx := s.auditDetailRepo.(*repository.InventoryAuditDetailRepository).WithTx(tx)

		// Get audit details
		details, err := auditDetailRepoTx.GetByAuditId(int(audit.AuditID))
		if err != nil {
			tx.Rollback()
			return common.SystemError
		}

		binRepo := repository.NewBinRepository()
		userID, _ := strconv.Atoi(middleware.GetUserID(c))

		// Update component bins and create ledger entries
		for _, detail := range details {
			compBin, err := compBinRepoTx.GetByComponentAndBinId(detail.ComponentID, detail.BinID)
			if err != nil {
				tx.Rollback()
				return common.SystemError
			}

			if compBin == nil {
				// Create if not exists
				compBin = &model.ComponentBin{
					ComponentID: detail.ComponentID,
					BinID:       detail.BinID,
					Quantity:    detail.ActualQuantity,
					CreatedAt:   time.Now(),
					CreatedBy:   userID,
				}
				err = compBinRepoTx.Save(compBin)
			} else {
				compBin.Quantity = detail.ActualQuantity
				compBin.UpdatedAt = time.Now()
				compBin.UpdatedBy = userID
				err = compBinRepoTx.Update(compBin)
			}

			if err != nil {
				tx.Rollback()
				return common.SystemError
			}

			// Create ledger entry if there's a difference
			if detail.DifferenceQuantity != 0 {
				binInfo, _ := binRepo.GetById(detail.BinID)
				warehouseID := 0
				if binInfo != nil {
					warehouseID = binInfo.WarehouseID
				}

				ledgerReq := &dtos.InventoryLedgerCreate{
					ComponentID:     detail.ComponentID,
					WarehouseID:     warehouseID,
					BinID:           detail.BinID,
					ReferenceType:   constants.LedgerReferenceTypeAudit,
					ReferenceTypeID: audit.AuditID,
					Description:     fmt.Sprintf("Kiểm kê số #%d", audit.AuditID),
					QuantityChange:  detail.DifferenceQuantity,
					QuantityAfter:   detail.ActualQuantity,
					Note:            detail.Note,
					CreatedBy:       userID,
				}
				if err := s.ledgerService.CreateInventoryLedgerEntry(ledgerReq); err != nil {
					tx.Rollback()
					return &common.Error{Code: "500", Message: fmt.Sprintf("Lỗi tạo sổ cái: %v", err)}
				}
			}
		}

		// Update audit status
		audit.StatusInt = req.StatusInt
		audit.UpdatedBy = userID
		audit.UpdatedAt = time.Now()
		err = s.auditRepo.Update(audit)
		if err != nil {
			tx.Rollback()
			return common.SystemError
		}

		if err := tx.Commit().Error; err != nil {
			return common.SystemError
		}
	} else {
		// For other status changes, just update the audit without ledger entries
		userID, _ := strconv.Atoi(middleware.GetUserID(c))
		audit.StatusInt = req.StatusInt
		audit.UpdatedBy = userID
		audit.UpdatedAt = time.Now()

		err = s.auditRepo.Update(audit)
		if err != nil {
			return common.SystemError
		}
	}

	return nil
}

func (s *InventoryAuditService) CreateInventoryAuditDetail(c *gin.Context) *common.Error {
	var req []dtos.InventoryAuditDetailCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if len(req) == 0 {
		return common.RequestInvalid
	}

	// Check audit exists
	_, err := s.auditRepo.GetByRequestId(req[0].AuditID)
	if err != nil {
		return common.NotFound
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	details := make([]model.InventoryAuditDetail, 0)
	for _, r := range req {
		auditDetail := model.InventoryAuditDetail{
			AuditID:            r.AuditID,
			ComponentID:        r.ComponentID,
			BinID:              r.BinID,
			SystemQuantity:     r.SystemQuantity,
			ActualQuantity:     r.ActualQuantity,
			DifferenceQuantity: math.Abs(r.SystemQuantity - r.ActualQuantity),
			CreatedBy:          userID,
		}
		details = append(details, auditDetail)
	}

	if err := s.auditRepo.CreateInventoryAuditDetailsTx(details); err != nil {
		return common.SystemError
	}

	return nil
}
func (s *InventoryAuditService) GetAllInventoryAuditDetails(c *gin.Context) ([]dtos.InventoryAuditDetailResponse, int, *common.Error) {
	var query dtos.InventoryAuditDetailFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	auditDetails, total, err := s.auditDetailRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	var responses []dtos.InventoryAuditDetailResponse
	for _, auditDetail := range auditDetails {
		responses = append(responses, modelToInventoryAuditDetailResponse(&auditDetail))
	}
	return responses, total, nil
}

func (s *InventoryAuditService) UpdateInventoryAuditDetail(c *gin.Context) *common.Error {
	var req []dtos.InventoryAuditDetailUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	// Validate all items have AuditDetailID
	for _, r := range req {
		if r.AuditDetailID == 0 {
			return common.RequestInvalid
		}
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	detailRepoTx := s.auditDetailRepo.(*repository.InventoryAuditDetailRepository).WithTx(tx)

	// Update each audit detail
	for _, r := range req {
		auditDetail, err := detailRepoTx.GetByInventoryAuditDetailId(r.AuditDetailID)
		if err != nil {
			tx.Rollback()
			return common.NotFound
		}
		if auditDetail == nil {
			tx.Rollback()
			return &common.Error{Code: "404", Message: "Chi tiết kiểm kê không tồn tại"}
		}

		// Update fields if provided
		if r.ComponentID != 0 {
			auditDetail.ComponentID = r.ComponentID
		}
		if r.BinID != 0 {
			auditDetail.BinID = r.BinID
		}
		if r.SystemQuantity != 0 {
			auditDetail.SystemQuantity = r.SystemQuantity
		}
		if r.ActualQuantity != 0 {
			auditDetail.ActualQuantity = r.ActualQuantity
		}
		// Recalculate difference if quantities are provided
		if r.SystemQuantity != 0 && r.ActualQuantity != 0 {
			auditDetail.DifferenceQuantity = math.Abs(float64(r.SystemQuantity - r.ActualQuantity))
		}
		if r.UpdatedBy != 0 {
			auditDetail.UpdatedBy = r.UpdatedBy
		}
		userID, _ := strconv.Atoi(middleware.GetUserID(c))
		auditDetail.UpdatedBy = userID
		auditDetail.UpdatedAt = time.Now()

		err = detailRepoTx.Update(auditDetail)
		if err != nil {
			tx.Rollback()
			return common.SystemError
		}
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func (s *InventoryAuditService) DeleteInventoryAuditDetail(c *gin.Context) *common.Error {
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

	err := s.auditDetailRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}
func modelToInventoryAuditDetailResponse(auditDetail *model.InventoryAuditDetail) dtos.InventoryAuditDetailResponse {
	return dtos.InventoryAuditDetailResponse{
		AuditDetailID:      auditDetail.AuditDetailID,
		AuditID:            auditDetail.AuditID,
		ComponentID:        auditDetail.ComponentID,
		BinID:              auditDetail.BinID,
		SystemQuantity:     auditDetail.SystemQuantity,
		ActualQuantity:     auditDetail.ActualQuantity,
		DifferenceQuantity: auditDetail.DifferenceQuantity,
		CreatedBy:          auditDetail.CreatedBy,
		CreatedAt:          auditDetail.CreatedAt,
		UpdatedBy:          auditDetail.UpdatedBy,
		UpdatedAt:          auditDetail.UpdatedAt,
	}
}
