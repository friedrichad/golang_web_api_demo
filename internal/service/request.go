package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/middleware"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/friedrichad/golang_web_api_demo/internal/redis"
)

type IRequestService interface {
	GetAllRequests(c *gin.Context) ([]model.RequestResponse, int, *common.Error)
	GetRequestById(c *gin.Context) (*model.RequestResponse, *common.Error)
	CreateRequest(c *gin.Context) (*model.RequestResponse, *common.Error)
	UpdateRequest(c *gin.Context) *common.Error
	DeleteRequest(c *gin.Context) *common.Error
	ApprovalRequest(c *gin.Context) *common.Error
	ConfirmRequest(c *gin.Context) *common.Error
	// RequestDetail CRUD operations
	GetAllRequestDetails(c *gin.Context) ([]model.RequestDetailResponse, int, *common.Error)
	GetRequestDetailById(c *gin.Context) (*model.RequestDetailResponse, *common.Error)
	CreateRequestDetail(c *gin.Context) (*model.RequestDetailResponse, *common.Error)
	UpdateRequestDetail(c *gin.Context) *common.Error
	DeleteRequestDetail(c *gin.Context) *common.Error
	ExpireRequests() error
}

type RequestService struct {
	requestRepo       repository.IRequestRepository
	requestDetailRepo repository.IRequestDetailRepository
	ledgerService     IInventoryLedgerService
	userService       IUserService
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

func (s *RequestService) GetAllRequests(c *gin.Context) ([]model.RequestResponse, int, *common.Error) {
	var query model.RequestFilter
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

	requestResponses := make([]model.RequestResponse, len(requests))
	for i, req := range requests {
		requestResponses[i] = modelToRequestResponse(&req)
	}

	return requestResponses, total, nil
}

func (s *RequestService) GetRequestById(c *gin.Context) (*model.RequestResponse, *common.Error) {
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

func (s *RequestService) CreateRequest(c *gin.Context) (*model.RequestResponse, *common.Error) {

	var req model.RequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))

	request := &model.Request{
		RequestType:   *req.RequestType,
		Description:   req.Description,
		WarehouseID:   int(req.WarehouseID),
		PerformedByID: userID,
		PartnerID:     int(req.PartnerID),
		ExpiredDate:   req.ExpiredDate,
		Note:          req.Note,
		StatusInt:     constants.RequestStatusPending, // Status = PENDING initially
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
	}

	err := s.requestRepo.Save(request)
	if err != nil {
		return nil, common.SystemError
	}
	requestResponse := modelToRequestResponse(request)
	return &requestResponse, nil
}

func (s *RequestService) UpdateRequest(c *gin.Context) *common.Error {
	var req model.RequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return &common.Error{Code: "400", Message: err.Error()}
	}

	request, err := s.requestRepo.GetByRequestId(int(req.RequestID))
	if err != nil {
		return common.NotFound
	}

	if request == nil {
		return common.NotFound
	}

	if req.RequestType != nil {
		request.RequestType = *req.RequestType
	}
	if req.Description != "" {
		request.Description = req.Description
	}
	if req.WarehouseID != 0 {
		request.WarehouseID = int(req.WarehouseID)
	}

	userID, _ := strconv.Atoi(middleware.GetUserID(c))

	if req.PartnerID != 0 {
		request.PartnerID = int(req.PartnerID)
	}
	if req.StatusInt != 0 {
		request.StatusInt = int(req.StatusInt)
	}
	if req.Note != "" {
		request.Note = req.Note
	}
	request.UpdatedBy = userID
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

func (s *RequestService) ApprovalRequest(c *gin.Context) *common.Error {
	var req model.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	request, err := s.requestRepo.GetByRequestId(int(req.RequestID))
	requesterInfo, err := s.userService.GetUserInfoWithCache(request.PerformedByID)
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin requester ID=%d: %v", request.PerformedByID, err)
		return common.SystemError
	}
	requesterPositionLevel := requesterInfo.PositionInfo.PositionLevel
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	check, err := redis.CanApproveRequest(redis.Rdb, userID, requesterPositionLevel)
	if err != nil {
		log.Print("Lỗi kiểm tra quyền phê duyệt: ", err)
		return common.RequestInvalid
	}
	if !check {
		log.Print("Người dùng không có quyền phê duyệt request này")
		return common.RequestInvalid
	}
	if err != nil {
		return common.SystemError
	}
	if request == nil {
		return common.NotFound
	}

	if request.StatusInt != constants.RequestStatusPending {
		return &common.Error{Code: "400", Message: "Chỉ yêu cầu PENDING mới được phê duyệt!"}
	}

	if req.StatusInt != constants.RequestStatusApproved && req.StatusInt != constants.RequestStatusRejected {
		return &common.Error{Code: "400", Message: "Yêu cầu chỉ có thể được phê duyệt (APPROVED) hoặc từ chối (REJECTED)!"}
	}

	request.ApproverID = userID
	request.StatusInt = req.StatusInt
	request.Reason = req.Reason
	request.UpdatedBy = userID
	request.UpdatedAt = time.Now()

	err = s.requestRepo.Update(request)
	if err != nil {
		return common.SystemError
	}
	return nil
}

func (s *RequestService) ConfirmRequest(c *gin.Context) *common.Error {
	var req model.ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if req.RequestID == 0 {
		return common.RequestInvalid
	}

	request, err := s.requestRepo.GetByRequestId(req.RequestID)
	if err != nil {
		return common.SystemError
	}
	if request == nil {
		return common.NotFound
	}

	if request.StatusInt != constants.RequestStatusApproved {
		return &common.Error{Code: "400", Message: "Chỉ yêu cầu APPROVED mới được xác nhận!"}
	}
	if req.StatusInt == constants.RequestStatusCompleted {
		if err := s.ApplyRequestDetails(c, req.RequestID); err != nil {
			return &common.Error{Code: "400", Message: err.Error()}
		}
	}
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	request.StatusInt = req.StatusInt
	request.UpdatedBy = userID
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
	userID, _ := strconv.Atoi(middleware.GetUserID(c))

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

			ledgerReq := &model.InventoryLedgerCreate{
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

			ledgerReq := &model.InventoryLedgerCreate{
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

func (s *RequestService) GetAllRequestDetails(c *gin.Context) ([]model.RequestDetailResponse, int, *common.Error) {
	var query model.RequestDetailFilter
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

	detailResponses := make([]model.RequestDetailResponse, len(details))
	for i, detail := range details {
		detailResponses[i] = modelToRequestDetailResponse(&detail)
	}

	return detailResponses, total, nil
}

func (s *RequestService) GetRequestDetailById(c *gin.Context) (*model.RequestDetailResponse, *common.Error) {
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

func (s *RequestService) CreateRequestDetail(c *gin.Context) (*model.RequestDetailResponse, *common.Error) {
	var req model.RequestDetailCreate
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

	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	detail := &model.RequestDetail{
		RequestID:   req.RequestID,
		ComponentID: req.ComponentID,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		BinFromID:   req.BinFromID,
		BinToID:     req.BinToID,
		CreatedBy:   userID,
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
	var req model.RequestDetailUpdate
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
	userID, _ := strconv.Atoi(middleware.GetUserID(c))
	detail.UpdatedBy = userID
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

func (s *RequestService) ExpireRequests() error {
	requests, err := s.requestRepo.GetExpiredPendingRequests()
	if err != nil {
		log.Print("Lỗi khi lấy request hết hạn: ", err)
		return err
	}

	if len(requests) == 0 {
		return nil
	}

	for _, r := range requests {
		r.StatusInt = constants.RequestStatusExpired
		r.UpdatedAt = time.Now()

		err := s.requestRepo.Update(&r)
		if err != nil {
			log.Print("Lỗi update request expired ID:", r.RequestID, err)
			continue
		}

		log.Print("Expired request ID:", r.RequestID)
	}

	return nil
}

func modelToRequestDetailResponse(detail *model.RequestDetail) model.RequestDetailResponse {
	return model.RequestDetailResponse{
		RequestDetailID: detail.RequestDetailID,
		RequestID:       detail.RequestID,
		ComponentID:     detail.ComponentID,
		Quantity:        detail.Quantity,
		UnitPrice:       detail.UnitPrice,
		BinFromID:       detail.BinFromID,
		BinToID:         detail.BinToID,
		CreatedBy:       detail.CreatedBy,
		CreatedAt:       detail.CreatedAt,
		UpdatedBy:       detail.UpdatedBy,
		UpdatedAt:       detail.UpdatedAt,
	}
}

func modelToRequestResponse(request *model.Request) model.RequestResponse {
	return model.RequestResponse{
		RequestID:     request.RequestID,
		RequestType:   request.RequestType,
		Description:   request.Description,
		WarehouseID:   request.WarehouseID,
		PerformedByID: request.PerformedByID,
		ApproverID:    request.ApproverID,
		PartnerID:     request.PartnerID,
		ExpiredDate:   request.ExpiredDate,
		StatusInt:     request.StatusInt,
		Note:          request.Note,
		CreatedBy:     request.CreatedBy,
		CreatedAt:     request.CreatedAt,
		UpdatedBy:     request.UpdatedBy,
		UpdatedAt:     request.UpdatedAt,
	}
}
