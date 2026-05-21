package service

import (
	"log"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/middleware"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IRequestPermissionService interface {
	GetAllPermissionByCondition(c *gin.Context) ([]model.RequestPermissionResponse, int, *common.Error)
	Save(c *gin.Context) *common.Error
	Update(c *gin.Context) *common.Error
	Delete(c *gin.Context) *common.Error
	Approval(c *gin.Context) *common.Error
}

type RequestPermissionService struct {
	repoRequestPermission repository.IRequestPermissionRepository
	repoRequest           repository.IRequestRepository
	repoUserPermission    repository.IUserPermissionRepository
	repoMenuPermission    repository.IMenuPermissionRepository
	userService           IUserService
}

var requestPermissionService IRequestPermissionService

func NewRequestPermissionService() IRequestPermissionService {
	if requestPermissionService == nil {
		requestPermissionService = &RequestPermissionService{
			repoRequestPermission: repository.NewRequestPermissionRepository(),
			repoRequest:           repository.NewRequestRepository(),
			repoUserPermission:    repository.NewUserPermissionRepository(),
			repoMenuPermission:    repository.NewMenuPermissionRepository(),
			userService:           NewUserService(),
		}
	}
	return requestPermissionService
}

func (s *RequestPermissionService) GetAllPermissionByCondition(c *gin.Context) ([]model.RequestPermissionResponse, int, *common.Error) {
	var query model.RequestPermissionFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Print("Lỗi khi bind query: ", err)
		return nil, 0, common.RequestInvalid
	}
	requestPermissions, total, err := s.repoRequestPermission.GetAllByCondition(query)
	if err != nil {
		log.Print("Lỗi khi lấy dữ liệu: ", err)
		return nil, 0, common.SystemError
	}
	requestPermissionResponses := make([]model.RequestPermissionResponse, len(requestPermissions))
	for i, _ := range requestPermissionResponses {
		requestPermissionResponses[i] = model.RequestPermissionResponse{
			RequestPermissionID: requestPermissions[i].RequestPermissionID,
			RequestID:           requestPermissions[i].RequestID,
			MenuPermission:      requestPermissions[i].MenuPermissionID,
			Reason:              requestPermissions[i].Reason,
			CreatedAt:           requestPermissions[i].CreatedAt,
		}
	}
	return requestPermissionResponses, total, nil
}

func (s *RequestPermissionService) Save(c *gin.Context) *common.Error {
	var requestPermissionCreate model.RequestPermissionCreate
	if err := c.ShouldBindJSON(&requestPermissionCreate); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}
	requestPermission := model.RequestPermission{
		RequestID:        requestPermissionCreate.RequestID,
		MenuPermissionID: requestPermissionCreate.MenuPermission,
		Reason:           requestPermissionCreate.Reason,
		CreatedAt:        time.Now(),
	}
	err := s.repoRequestPermission.Save(&requestPermission)
	if err != nil {
		log.Print("Lỗi khi lưu request permission: ", err)
		return common.SystemError
	}
	return nil
}

func (s *RequestPermissionService) Update(c *gin.Context) *common.Error {
	var requestPermissionUpdate model.RequestPermissionUpdate
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
	request.MenuPermissionID = requestPermissionUpdate.MenuPermission
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
	var req model.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}
	if req.StatusInt != constants.RequestStatusApproved &&
		req.StatusInt != constants.RequestStatusRejected &&
		req.StatusInt != constants.RequestStatusRevoked {
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
		request.StatusInt == constants.RequestStatusRejected ||
		request.StatusInt == constants.RequestStatusRevoked {
		log.Print("Request đã được xử lý trước đó")
		return common.RequestInvalid
	}
	userID, err := strconv.Atoi(middleware.GetUserID(c))
	if err != nil {
		log.Print("Lỗi parse userID: ", err)
		return common.RequestInvalid
	}

	requesterInfo, err := s.userService.GetUserInfoWithCache(request.PerformedByID)
	if err != nil {
		log.Printf("Lỗi khi lấy thông tin requester ID=%d: %v", request.PerformedByID, err)
		return common.SystemError
	}
	requesterPositionLevel := requesterInfo.PositionInfo.PositionLevel

	check, err := redis.CanApproveRequest(redis.Rdb, userID, requesterPositionLevel)
	if err != nil {
		log.Print("Lỗi kiểm tra quyền phê duyệt: ", err)
		return common.RequestInvalid
	}
	if !check {
		log.Print("Người dùng không có quyền phê duyệt request này")
		return common.RequestInvalid
	}
	switch req.StatusInt {
	case constants.RequestStatusApproved:
		if err := s.Confirm(c, req.RequestID, userID, request.PerformedByID, request.ExpiredDate); err != nil {
			log.Print("Lỗi confirm request: ", err)
			return common.SystemError
		}
		request.ApproverID = userID
		request.StatusInt = constants.RequestStatusApproved

	case constants.RequestStatusRejected:
		request.ApproverID = userID
		request.StatusInt = constants.RequestStatusRejected

	case constants.RequestStatusRevoked:
		if err := s.RevokePermission(c, req.RequestID, request.PerformedByID); err != nil {
			log.Print("Lỗi revoke permission: ", err)
			return common.SystemError
		}
		request.ApproverID = userID
		request.StatusInt = constants.RequestStatusRevoked
	}

	if err := s.repoRequest.Update(request); err != nil {
		log.Print("Lỗi update request: ", err)
		return common.SystemError
	}

	return nil
}

func (s *RequestPermissionService) Confirm(c *gin.Context, requestId int, approverId int, requesterId int, expiredDate time.Time) error {

	requestPermissions, err := s.repoRequestPermission.GetRequestPermissionByRequestId(requestId)
	if err != nil {
		log.Print("Có lỗi xảy ra khi truy vấn request: ", err)
		return err
	}

	if len(requestPermissions) == 0 {
		log.Print("Không tìm thấy request permissions cho request: ", requestId)
		return nil
	}
	userPermissions := make([]model.UserPermission, 0, len(requestPermissions))
	for _, rp := range requestPermissions {
		mp, err := s.repoMenuPermission.GetMenuPermissionById(rp.MenuPermissionID)
		if err != nil {
			log.Printf("Không tìm thấy MenuPermission với ID=%d: %v", rp.MenuPermissionID, err)
			continue
		}

		userPermissions = append(userPermissions, model.UserPermission{
			UserID:           requesterId,
			MenuPermissionID: mp.MenuPermissionID,
			ExpiredDate:      expiredDate,
			CreatedBy:        approverId,
			CreatedAt:        time.Now(),
			UpdatedBy:        approverId,
			UpdatedAt:        time.Now(),
		})
	}
	err = s.repoUserPermission.SaveBatch(userPermissions)
	if err != nil {
		log.Print("Có lỗi xảy ra khi lưu batch user permissions: ", err)
		return err
	}
	err = redis.DeleteUserPermissionField(redis.Rdb, requesterId)
	if err != nil {
		log.Print("Có lỗi xảy ra khi xóa user permission từ cache: ", err)
		return err
	}
	return nil
}
func (s *RequestPermissionService) RevokePermission(c *gin.Context, requestId int, requesterId int) error {
	requestPermissions, err := s.repoRequestPermission.GetRequestPermissionByRequestId(requestId)
	if err != nil {
		log.Print("Có lỗi xảy ra khi truy vấn request: ", err)
		return err
	}
	if len(requestPermissions) == 0 {
		log.Print("Không tìm thấy request permissions cho request: ", requestId)
		return nil
	}
	userPermissionIds := make([]int, 0, len(requestPermissions))
	for _, rp := range requestPermissions {
		mp, err := s.repoMenuPermission.GetMenuPermissionById(rp.MenuPermissionID)
		if err != nil {
			log.Printf("Không tìm thấy MenuPermission với ID=%d: %v", rp.MenuPermissionID, err)
			continue
		}
		userPermissionIds = append(userPermissionIds, mp.MenuPermissionID)
	}
	err = s.repoUserPermission.Delete(userPermissionIds)
	if err != nil {
		log.Print("Có lỗi xảy ra khi xóa user permissions: ", err)
		return err
	}
	err = redis.DeleteUserPermissionField(redis.Rdb, requesterId)
	if err != nil {
		log.Print("Có lỗi xảy ra khi xóa user permission từ cache: ", err)
		return err
	}
	return nil
}
