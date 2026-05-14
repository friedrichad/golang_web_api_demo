package service

import (
	"log"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/middleware"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IRequestPermissionService interface {
	GetAllPermissionByCondition(c *gin.Context) ([]dtos.RequestPermissionResponse, int, *common.Error)
	Save(c *gin.Context) *common.Error
	Update(c *gin.Context) *common.Error
	Delete(c *gin.Context) *common.Error
	Approval(c *gin.Context) *common.Error
}

type RequestPermissionService struct {
	repoRequestPermission repository.IRequestPermissionRepository
	repoRequest repository.IRequestRepository
	repoUserPermission repository.IUserPermissionRepository
}

var requestPermissionService IRequestPermissionService

func NewRequestPermissionService() IRequestPermissionService {
	if requestPermissionService == nil {
		requestPermissionService = &RequestPermissionService{
			repoRequestPermission: repository.NewRequestPermissionRepository(),
			repoRequest:                  repository.NewRequestRepository(),
			repoUserPermission: repository.NewUserPermissionRepository(),
		}
	}
	return requestPermissionService
}

func (s *RequestPermissionService) GetAllPermissionByCondition(c *gin.Context) ([]dtos.RequestPermissionResponse, int, *common.Error) {
	var query dtos.RequestPermissionFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Print("Lỗi khi bind query: ", err)
		return nil, 0, common.RequestInvalid
	}
	requestPermissions, total, err := s.repoRequestPermission.GetAllByCondition(query)
	if err != nil {
		log.Print("Lỗi khi lấy dữ liệu: ", err)
		return nil, 0, common.SystemError
	}
	requestPermissionResponses := make([]dtos.RequestPermissionResponse, len(requestPermissions))
	for i, _ := range requestPermissionResponses {
		requestPermissionResponses[i] = dtos.RequestPermissionResponse{
			RequestPermissionID: requestPermissions[i].RequestPermissionID,
			RequestID:           requestPermissions[i].RequestID,
			MenuID:              requestPermissions[i].MenuID,
			PermissionID:        requestPermissions[i].PermissionID,
			Reason:              requestPermissions[i].Reason,
			CreatedAt:           requestPermissions[i].CreatedAt,
		}
	}
	return requestPermissionResponses, total, nil
}

func (s *RequestPermissionService) Save(c *gin.Context) *common.Error {
	var requestPermissionCreate dtos.RequestPermissionCreate
	if err := c.ShouldBindJSON(&requestPermissionCreate); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}
	requestPermission := model.RequestPermission{
		RequestID:    requestPermissionCreate.RequestID,
		MenuID:       requestPermissionCreate.MenuID,
		PermissionID: requestPermissionCreate.PermissionID,
		Reason:       requestPermissionCreate.Reason,
		CreatedAt:    time.Now(),
	}
	err := s.repoRequestPermission.Save(&requestPermission)
	if err != nil {
		log.Print("Lỗi khi lưu request permission: ", err)
		return common.SystemError
	}
	return nil
}

func (s *RequestPermissionService) Update(c *gin.Context) *common.Error {
	var requestPermissionUpdate dtos.RequestPermissionUpdate
	if err := c.ShouldBindJSON(&requestPermissionUpdate); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}

	var request *model.RequestPermission
	request, err := s.repoRequestPermission.GetByRequestPermissionId(requestPermissionUpdate.RequestPermissionID)
	if err == gorm.ErrRecordNotFound {
		log.Print("Không tìm thấy request permission: ", err)
		return common.NotFound
	}
	if err != nil {
		log.Print("Lỗi khi lấy request permission: ", err)
		return common.SystemError
	}
	request.MenuID = requestPermissionUpdate.MenuID
	request.PermissionID = requestPermissionUpdate.PermissionID
	request.Reason = requestPermissionUpdate.Reason
	err = s.repoRequestPermission.Update(request)
	if err != nil {
		log.Print("Lỗi khi cập nhật request permission: ", err)
		return common.SystemError
	}
	return nil
}

func (s *RequestPermissionService) Delete(c *gin.Context) *common.Error {
	var ids struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&ids); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}
	err := s.repoRequestPermission.Delete(ids.IDs)
	if err != nil {
		log.Print("Lỗi khi xóa request permission: ", err)
		return common.SystemError
	}
	return nil
}

func (s *RequestPermissionService) Approval(c *gin.Context) *common.Error {
	var req dtos.ApprovalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}

	if req.StatusInt != constants.RequestStatusApproved &&
		req.StatusInt != constants.RequestStatusRejected {
		log.Print("Trạng thái phê duyệt không hợp lệ: ", req.StatusInt)
		return common.RequestInvalid
	}

	if req.Reason == "" {
		log.Print("Lý do phê duyệt không được để trống")
		return common.RequestInvalid
	}

	request, err := s.repoRequest.GetByRequestId(req.RequestID)
	if err != nil {
		log.Print("Không tìm thấy request: ", err)
		return common.NotFound
	}

	if request.StatusInt == constants.RequestStatusApproved ||
		request.StatusInt == constants.RequestStatusRejected {
		log.Print("Request đã được xử lý trước đó")
		return common.RequestInvalid
	}

	userID, err := strconv.Atoi(middleware.GetUserID(c))
	if err != nil {
		log.Print("Lỗi parse userID: ", err)
		return common.RequestInvalid
	}

	ok, err := s.repoRequest.CanApprove(userID, request.PerformedByID)
	if err != nil {
		log.Print("Lỗi kiểm tra quyền phê duyệt: ", err)
		return common.SystemError
	}

	if !ok {
		log.Print("User không có quyền phê duyệt request")
		return common.UserForbidden
	}

	switch req.StatusInt {

	case constants.RequestStatusApproved:
		if err := s.Confirm(c, req.RequestID, userID, request.PerformedByID); err != nil {
			log.Print("Lỗi confirm request: ", err)
			return common.SystemError
		}

		request.StatusInt = constants.RequestStatusApproved

	case constants.RequestStatusRejected:
		request.StatusInt = constants.RequestStatusRejected
	}

	if err := s.repoRequest.Update(request); err != nil {
		log.Print("Lỗi update request: ", err)
		return common.SystemError
	}

	return nil
}

func (s *RequestPermissionService) Confirm(c *gin.Context, requestId int, approverId int, requesterId int) error {
	requestPermissions, err := s.repoRequestPermission.GetRequestPermissionByRequestId(requestId)
	if err != nil {
		log.Print("Có lỗi xảy ra khi truy vấn request: ", err)
		return err
	}
	for _, rp := range requestPermissions {
		var req = model.UserPermission{
			UserID:    requesterId,
			MenuID:    rp.MenuID,
			PermissionID: rp.PermissionID,
			UpdatedBy: approverId,
			UpdatedAt: time.Now(),
		}
		err := s.repoUserPermission.Save(&req)
		if err != nil {
			log.Print("Có lỗi xảy ra khi lưu user permission: Menu ", rp.MenuID, ", Permission ", rp.PermissionID, ": ", err)
			return err
		}
	} 
	return nil
}