package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IRequestService interface {
	GetAllRequests(c *gin.Context) ([]dtos.RequestResponse, int, *common.Error)
	GetRequestById(c *gin.Context) (*dtos.RequestResponse, *common.Error)
	CreateRequest(c *gin.Context) (*dtos.RequestResponse, *common.Error)
	UpdateRequest(c *gin.Context) *common.Error
	DeleteRequest(c *gin.Context) *common.Error
	ApprovalRequest(c *gin.Context) *common.Error
	ConfirmRequest(c *gin.Context) *common.Error
	// RequestDetail CRUD operations
	GetAllRequestDetails(c *gin.Context) ([]dtos.RequestDetailResponse, int, *common.Error)
	GetRequestDetailById(c *gin.Context) (*dtos.RequestDetailResponse, *common.Error)
	CreateRequestDetail(c *gin.Context) (*dtos.RequestDetailResponse, *common.Error)
	UpdateRequestDetail(c *gin.Context) *common.Error
	DeleteRequestDetail(c *gin.Context) *common.Error
}

type RequestService struct {
	requestRepo       repository.IRequestRepository
	requestDetailRepo repository.IRequestDetailRepository
	ledgerService     IInventoryLedgerService
}

var requestService IRequestService

func NewRequestService() IRequestService {
	if requestService == nil {
		requestService = &RequestService{
			requestRepo:       repository.NewRequestRepository(),
			requestDetailRepo: repository.NewRequestDetailRepository(),
			ledgerService:     NewInventoryLedgerService(),
		}
	}
	return requestService
}

func (s *RequestService) GetAllRequests(c *gin.Context) ([]dtos.RequestResponse, int, *common.Error) {
	var query dtos.RequestFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	requests, total, err := s.requestRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	requestResponses := make([]dtos.RequestResponse, len(requests))
	for i, req := range requests {
		requestResponses[i] = modelToRequestResponse(&req)
	}

	return requestResponses, total, nil
}

func (s *RequestService) GetRequestById(c *gin.Context) (*dtos.RequestResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	requestId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	request, err := s.requestRepo.GetByRequestId(int(requestId))
	if err != nil {
		return nil, common.NotFound
	}

	if request == nil {
		return nil, &common.Error{Code: "404", Message: "Yêu cầu không tồn tại"}
	}

	requestResponse := modelToRequestResponse(request)
	return &requestResponse, nil
}

func (s *RequestService) CreateRequest(c *gin.Context) (*dtos.RequestResponse, *common.Error) {
	var req dtos.RequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	requestRepoTx := s.requestRepo.(*repository.RequestRepository).WithTx(tx)
	requestDetailRepoTx := s.requestDetailRepo.(*repository.RequestDetailRepository).WithTx(tx)

	request := &model.Request{
		RequestType:   req.RequestType,
		Description:   req.Description,
		WarehouseID:   int(req.WarehouseID),
		PerformedByID: int(req.PerformedByID),
		PartnerID:     int(req.PartnerID),
		RequestDate:   time.Now(),
		Note:          req.Note,
		StatusInt:     1,
		CreatedAt:     time.Now(),
	}

	err := requestRepoTx.Save(request)
	if err != nil {
		tx.Rollback()
		return nil, common.SystemError
	}
	for _, detailDto := range req.RequestDetail {
		if detailDto.ComponentID == nil || detailDto.Quantity == nil || detailDto.UnitPrice == nil {
			tx.Rollback()
			return nil, common.RequestInvalid
		}

		detail := &model.RequestDetail{
			RequestID:   int(request.RequestID),
			ComponentID: *detailDto.ComponentID,
			Quantity:    *detailDto.Quantity,
			UnitPrice:   *detailDto.UnitPrice,
		}
		if detailDto.BinFromID != nil {
			detail.BinFromID = *detailDto.BinFromID
		}
		if detailDto.BinToID != nil {
			detail.BinToID = *detailDto.BinToID
		}

		err = requestDetailRepoTx.Save(detail)
		if err != nil {
			tx.Rollback()
			return nil, common.SystemError
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, common.SystemError
	}

	requestResponse := modelToRequestResponse(request)
	return &requestResponse, nil
}

func (s *RequestService) UpdateRequest(c *gin.Context) *common.Error {
	var req dtos.RequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	request, err := s.requestRepo.GetByRequestId(int(req.RequestID))
	if err != nil {
		return common.NotFound
	}

	if request == nil {
		return common.NotFound
	}

	if req.RequestType != "" {
		request.RequestType = req.RequestType
	}
	if req.Description != "" {
		request.Description = req.Description
	}
	if req.WarehouseID != 0 {
		request.WarehouseID = int(req.WarehouseID)
	}
	if req.PerformedByID != 0 {
		request.PerformedByID = int(req.PerformedByID)
	}
	if req.ApproverID != 0 {
		request.ApproverID = int(req.ApproverID)
	}
	if req.PartnerID != 0 {
		request.PartnerID = int(req.PartnerID)
	}
	if req.StatusInt != 0 {
		request.StatusInt = int(req.StatusInt)
	}
	if req.Note != "" {
		request.Note = req.Note
	}
	request.UpdatedBy = int(req.UpdatedBy)
	request.UpdatedAt = time.Now()

	err = s.requestRepo.Update(request)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *RequestService) DeleteRequest(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		requestId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(requestId)
	}

	err := s.requestRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func modelToRequestResponse(request *model.Request) dtos.RequestResponse {
	return dtos.RequestResponse{
		RequestID:     request.RequestID,
		RequestType:   request.RequestType,
		Description:   request.Description,
		WarehouseID:   request.WarehouseID,
		PerformedByID: request.PerformedByID,
		ApproverID:    request.ApproverID,
		PartnerID:     request.PartnerID,
		RequestDate:   request.RequestDate,
		StatusInt:     request.StatusInt,
		Note:          request.Note,
		CreatedAt:     request.CreatedAt,
		CreateBy:      request.CreatedBy,
		UpdatedAt:     request.UpdatedAt,
		UpdatedBy:     request.UpdatedBy,
	}
}
func (s *RequestService) ApprovalRequest(c *gin.Context) *common.Error {
	var req dtos.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}
	request, err := s.requestRepo.GetByRequestId(int(req.RequestID))
	if err != nil {
		return common.SystemError
	}
	if request == nil {
		return common.NotFound
	}
	request.ApproverID = int(req.ApproverID)
	if !constants.IsValidRequestStatus(req.StatusInt) {
		return &common.Error{Code: "400", Message: "Trạng thái yêu cầu không hợp lệ!"}
	}
	if request.StatusInt == constants.RequestStatusApproved || request.StatusInt == constants.RequestStatusRejected {
		return &common.Error{Code: "400", Message: "Không thể phê duyệt lại yêu cầu đã được phê duyệt!"}
	}
	if req.StatusInt != constants.RequestStatusApproved && req.StatusInt != constants.RequestStatusRejected {
		return &common.Error{Code: "400", Message: "Yêu cầu chỉ có thể được phê duyệt hoặc từ chối!"}
	}
	request.StatusInt = int(req.StatusInt)
	request.Note = req.Note
	request.UpdatedAt = time.Now()
	err = s.requestRepo.Update(request)
	if err != nil {
		return common.SystemError
	}
	return nil
}

func (s *RequestService) ConfirmRequest(c *gin.Context) *common.Error {
	var req dtos.ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}
	var request *model.Request
	if req.RequestID == 0 {
		return common.RequestInvalid
	}
	request, err := s.requestRepo.GetByRequestId(int(req.RequestID))
	if err != nil {
		return common.SystemError
	}
	if request == nil {
		return common.NotFound
	}
	if request.StatusInt != constants.RequestStatusCompleted && request.StatusInt != constants.RequestStatusCancelled {
		return &common.Error{Code: "400", Message: "Đơn đã hoàn thành/hủy, không thể chỉnh sửa!"}
	}
	if !constants.IsValidRequestStatus(req.StatusInt) {
		return &common.Error{Code: "400", Message: "Trạng thái yêu cầu không hợp lệ!"}
	}

	if request.StatusInt != constants.RequestStatusApproved {
		return &common.Error{Code: "400", Message: "Chỉ có thể xác nhận yêu cầu đã được phê duyệt!"}
	}
	if req.StatusInt == constants.RequestStatusCompleted {
		if err := s.ApplyRequestDetails(c, int(req.RequestID)); err != nil {
			return &common.Error{Code: "400", Message: err.Error()}
		}
	}
	request.StatusInt = req.StatusInt
	request.UpdatedAt = time.Now()
	err = s.requestRepo.Update(request)
	if err != nil {
		return common.SystemError
	}
	return nil
}
func (s *RequestService) ApplyRequestDetails(c *gin.Context, requestID int) error {
	details, err := s.requestDetailRepo.GetByRequestId(requestID)
	if err != nil {
		return err
	}

	// Get request to know the type
	request, err := s.requestRepo.GetByRequestId(requestID)
	if err != nil || request == nil {
		return fmt.Errorf("yêu cầu không tồn tại")
	}

	// Get user ID from context for ledger
	userIDInterface, _ := c.Get("user_id")
	userID := 0
	if userIDInterface != nil {
		userID = userIDInterface.(int)
	}

	compBinRepo := repository.NewComponentBinRepository()
	binRepo := repository.NewBinRepository()

	for _, detail := range details {

		if detail.BinFromID > 0 {
			fromBin, err := compBinRepo.GetByComponentAndBinId(detail.ComponentID, detail.BinFromID)
			if err != nil {
				return err
			}

			if fromBin == nil {
				return fmt.Errorf("bin %d không tồn tại", detail.BinFromID)
			}

			if fromBin.Quantity < float64(detail.Quantity) {
				return fmt.Errorf("không đủ hàng trong bin %d, hiện có %.0f, cần %d", detail.BinFromID, fromBin.Quantity, detail.Quantity)
			}

			fromBin.Quantity -= float64(detail.Quantity)
			if err := compBinRepo.Update(fromBin); err != nil {
				return err
			}

			// Create ledger entry for source bin
			binInfo, _ := binRepo.GetById(detail.BinFromID)
			warehouseID := 0
			if binInfo != nil {
				warehouseID = binInfo.WarehouseID
			}

			ledgerReq := &dtos.InventoryLedgerCreate{
				ComponentID:     detail.ComponentID,
				WarehouseID:     warehouseID,
				BinID:           detail.BinFromID,
				ReferenceType:   constants.LedgerReferenceTypeRequest,
				ReferenceTypeID: requestID,
				Description:     fmt.Sprintf("Export từ yêu cầu #%d", requestID),
				QuantityChange:  -float64(detail.Quantity),
				QuantityAfter:   fromBin.Quantity,
				Note:            request.Note,
				CreatedBy:       userID,
			}
			if err := s.ledgerService.CreateInventoryLedgerEntry(ledgerReq); err != nil {
				return fmt.Errorf("lỗi tạo sổ cái: %v", err)
			}
		}

		if detail.BinToID > 0 {
			toBin, err := compBinRepo.GetByComponentAndBinId(detail.ComponentID, detail.BinToID)
			if err != nil {
				return err
			}

			newQuantity := float64(detail.Quantity)
			if toBin == nil {
				newBin := &model.ComponentBin{
					ComponentID: detail.ComponentID,
					BinID:       detail.BinToID,
					Quantity:    newQuantity,
					CreatedBy:   detail.CreatedBy,
					CreatedAt:   time.Now(),
				}
				if err := compBinRepo.Save(newBin); err != nil {
					return err
				}
				toBin = newBin
			} else {
				toBin.Quantity += newQuantity
				if err := compBinRepo.Update(toBin); err != nil {
					return err
				}
			}

			// Create ledger entry for target bin
			binInfo, _ := binRepo.GetById(detail.BinToID)
			warehouseID := 0
			if binInfo != nil {
				warehouseID = binInfo.WarehouseID
			}

			ledgerReq := &dtos.InventoryLedgerCreate{
				ComponentID:     detail.ComponentID,
				WarehouseID:     warehouseID,
				BinID:           detail.BinToID,
				ReferenceType:   constants.LedgerReferenceTypeRequest,
				ReferenceTypeID: requestID,
				Description:     fmt.Sprintf("Import từ yêu cầu #%d", requestID),
				QuantityChange:  newQuantity,
				QuantityAfter:   toBin.Quantity,
				Note:            request.Note,
				CreatedBy:       userID,
			}
			if err := s.ledgerService.CreateInventoryLedgerEntry(ledgerReq); err != nil {
				return fmt.Errorf("lỗi tạo sổ cái: %v", err)
			}
		}
	}

	return nil
}

// RequestDetail CRUD operations

func (s *RequestService) GetAllRequestDetails(c *gin.Context) ([]dtos.RequestDetailResponse, int, *common.Error) {
	var query dtos.RequestDetailFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	details, total, err := s.requestDetailRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	detailResponses := make([]dtos.RequestDetailResponse, len(details))
	for i, detail := range details {
		detailResponses[i] = modelToRequestDetailResponse(&detail)
	}

	return detailResponses, total, nil
}

func (s *RequestService) GetRequestDetailById(c *gin.Context) (*dtos.RequestDetailResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	detailId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	detail, err := s.requestDetailRepo.GetByRequestDetailId(int(detailId))
	if err != nil {
		return nil, common.NotFound
	}

	if detail == nil {
		return nil, &common.Error{Code: "404", Message: "Chi tiết yêu cầu không tồn tại"}
	}

	detailResponse := modelToRequestDetailResponse(detail)
	return &detailResponse, nil
}

func (s *RequestService) CreateRequestDetail(c *gin.Context) (*dtos.RequestDetailResponse, *common.Error) {
	var req dtos.RequestDetailCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	// Verify request exists
	request, err := s.requestRepo.GetByRequestId(req.RequestID)
	if err != nil || request == nil {
		return nil, &common.Error{Code: "404", Message: "Yêu cầu không tồn tại"}
	}

	detail := &model.RequestDetail{
		RequestID:   req.RequestID,
		ComponentID: req.ComponentID,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		BinFromID:   req.BinFromID,
		BinToID:     req.BinToID,
		CreatedAt:   time.Now(),
	}

	err = s.requestDetailRepo.Save(detail)
	if err != nil {
		return nil, common.SystemError
	}

	detailResponse := modelToRequestDetailResponse(detail)
	return &detailResponse, nil
}

func (s *RequestService) UpdateRequestDetail(c *gin.Context) *common.Error {
	var req dtos.RequestDetailUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return &common.Error{Code: "400", Message: err.Error()}
	}

	detail, err := s.requestDetailRepo.GetByRequestDetailId(req.RequestDetailID)
	if err != nil || detail == nil {
		return &common.Error{Code: "404", Message: "Chi tiết yêu cầu không tồn tại"}
	}

	// Update fields if provided
	if req.RequestID != 0 {
		detail.RequestID = req.RequestID
	}
	if req.ComponentID != 0 {
		detail.ComponentID = req.ComponentID
	}
	if req.Quantity != 0 {
		detail.Quantity = req.Quantity
	}
	if req.UnitPrice != 0 {
		detail.UnitPrice = req.UnitPrice
	}
	if req.BinFromID != 0 {
		detail.BinFromID = req.BinFromID
	}
	if req.BinToID != 0 {
		detail.BinToID = req.BinToID
	}
	if req.UpdatedBy != 0 {
		detail.UpdatedBy = req.UpdatedBy
	}
	detail.UpdatedAt = time.Now()

	err = s.requestDetailRepo.Update(detail)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *RequestService) DeleteRequestDetail(c *gin.Context) *common.Error {
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

	err := s.requestDetailRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func modelToRequestDetailResponse(detail *model.RequestDetail) dtos.RequestDetailResponse {
	return dtos.RequestDetailResponse{
		RequestDetailID: detail.RequestDetailID,
		RequestID:       detail.RequestID,
		ComponentID:     detail.ComponentID,
		Quantity:        detail.Quantity,
		UnitPrice:       detail.UnitPrice,
		BinFromID:       detail.BinFromID,
		BinToID:         detail.BinToID,
	}
}
