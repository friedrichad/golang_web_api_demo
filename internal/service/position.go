package service

import (
	"log"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IPositionService interface {
	GetAllPositions(c *gin.Context) ([]model.PositionResponse, int, *common.Error)
	GetPositionById(c *gin.Context) (*model.PositionResponse, *common.Error)
	CreatePosition(c *gin.Context) (*model.PositionResponse, *common.Error)
	UpdatePosition(c *gin.Context) *common.Error
	DeletePosition(c *gin.Context) *common.Error
}

type PositionService struct {
	positionRepo repository.IPositionRepository
}

var positionService IPositionService

func NewPositionService() IPositionService {
	if positionService == nil {
		positionService = &PositionService{
			positionRepo: repository.NewPositionRepository(),
		}
	}
	return positionService
}

func (s *PositionService) GetAllPositions(c *gin.Context) ([]model.PositionResponse, int, *common.Error) {
	var query model.PositionFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Print("Lỗi khi bind query: ", err)
		return nil, 0, common.RequestInvalid
	}

	positions, total, err := s.positionRepo.GetAllByCondition(query)
	if err != nil {
		log.Print("Lỗi khi lấy dữ liệu: ", err)
		return nil, 0, common.SystemError
	}

	if total == 0 {
		return nil, 0, common.NotFound
	}

	positionResponses := make([]model.PositionResponse, len(positions))
	for i, position := range positions {
		positionResponses[i] = modelToPositionResponse(&position)
	}

	return positionResponses, total, nil
}

func (s *PositionService) GetPositionById(c *gin.Context) (*model.PositionResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	positionId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	position, err := s.positionRepo.GetPositionById(int(positionId))
	if err != nil {
		log.Print("Lỗi khi lấy position: ", err)
		return nil, common.NotFound
	}

	if position == nil {
		return nil, &common.Error{Code: "404", Message: "Vị trí không tồn tại"}
	}

	positionResponse := modelToPositionResponse(position)
	return &positionResponse, nil
}

func (s *PositionService) CreatePosition(c *gin.Context) (*model.PositionResponse, *common.Error) {
	var req model.PositionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return nil, common.RequestInvalid
	}

	position := &model.Position{
		PositionName: req.PositionName,
		Description:  req.Description,
		CreatedAt:    time.Now(),
	}

	err := s.positionRepo.Save(position)
	if err != nil {
		log.Print("Lỗi khi tạo position: ", err)
		return nil, common.SystemError
	}

	positionResponse := modelToPositionResponse(position)
	return &positionResponse, nil
}

func (s *PositionService) UpdatePosition(c *gin.Context) *common.Error {
	var req model.PositionUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}

	position, err := s.positionRepo.GetPositionById(req.PositionID)
	if err != nil {
		log.Print("Lỗi khi lấy position: ", err)
		return common.NotFound
	}

	if position == nil {
		return &common.Error{Code: "404", Message: "Vị trí không tồn tại"}
	}

	position.PositionName = req.PositionName
	position.Description = req.Description
	position.UpdatedAt = time.Now()

	err = s.positionRepo.Update(position)
	if err != nil {
		log.Print("Lỗi khi cập nhật position: ", err)
		return common.SystemError
	}

	return nil
}

func (s *PositionService) DeletePosition(c *gin.Context) *common.Error {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print("Lỗi khi bind json: ", err)
		return common.RequestInvalid
	}

	if len(req.IDs) == 0 {
		return common.RequestInvalid
	}

	err := s.positionRepo.Delete(req.IDs)
	if err != nil {
		log.Print("Lỗi khi xóa position: ", err)
		return common.SystemError
	}

	return nil
}

func modelToPositionResponse(position *model.Position) model.PositionResponse {
	return model.PositionResponse{
		PositionID:   position.PositionID,
		PositionName: position.PositionName,
		Description:  position.Description,
		CreatedBy:    position.CreatedBy,
		CreatedAt:    position.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:    position.UpdatedBy,
		UpdatedAt:    position.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
