package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IRequestService interface {
	GetAllRequests(c *gin.Context) ([]dtos.RequestResponse, int, *common.Error)
	GetRequestById(c *gin.Context) (*dtos.RequestResponse, *common.Error)
	CreateRequest(c *gin.Context) (*dtos.RequestResponse, *common.Error)
	UpdateRequest(c *gin.Context) *common.Error
	DeleteRequest(c *gin.Context) *common.Error
}

type RequestService struct {
	requestRepo       repository.IRequestRepository
	requestDetailRepo repository.IRequestDetailRepository
}

var requestService IRequestService

func NewRequestService() IRequestService {
	if requestService == nil {
		requestService = &RequestService{
			requestRepo:       repository.NewRequestRepository(),
			requestDetailRepo: repository.NewRequestDetailRepository(),
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

	requestId, err := strconv.ParseInt(idStr, 10, 32)
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
		BinFrom:       int(req.BinFrom),
		BinTo:         int(req.BinTo),
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
		return &common.Error{Code: "404", Message: "Yêu cầu không tồn tại"}
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
	if req.BinFrom != 0 {
		request.BinFrom = int(req.BinFrom)
	}
	if req.BinTo != 0 {
		request.BinTo = int(req.BinTo)
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
		requestId, err := strconv.ParseInt(idStr, 10, 32)
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
		RequestID:     int(request.RequestID),
		RequestType:   request.RequestType,
		Description:   request.Description,
		WarehouseID:   int(request.WarehouseID),
		BinTo:         int(request.BinTo),
		BinFrom:       int(request.BinFrom),
		PerformedByID: int(request.PerformedByID),
		ApproverID:    int(request.ApproverID),
		PartnerID:     int(request.PartnerID),
		RequestDate:   request.RequestDate,
		StatusInt:     int(request.StatusInt),
		Note:          request.Note,
		CreatedAt:     request.CreatedAt,
		CreateBy:      int(request.CreatedBy),
		UpdatedAt:     request.UpdatedAt,
		UpdatedBy:     int(request.UpdatedBy),
	}
}
