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

type IBinService interface {
	GetAllBins(c *gin.Context) ([]dtos.BinResponse, int, *common.Error)
	GetBinByBinId(c *gin.Context) (*dtos.BinResponse, *common.Error)
	CreateBin(c *gin.Context) (*dtos.BinResponse, *common.Error)
	UpdateBin(c *gin.Context) *common.Error
	DeleteBin(c *gin.Context) *common.Error
}

type BinService struct {
	binRepo repository.IBinRepository
}

var binService IBinService

func NewBinService() IBinService {
	if binService == nil {
		binService = &BinService{
			binRepo: repository.NewBinRepository(),
		}
	}
	return binService
}

func (s *BinService) GetAllBins(c *gin.Context) ([]dtos.BinResponse, int, *common.Error) {
	var query dtos.BinFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	bins, total, err := s.binRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	binResponses := make([]dtos.BinResponse, len(bins))
	for i, bin := range bins {
		binResponses[i] = modelToBinResponse(&bin)
	}

	return binResponses, total, nil
}

func (s *BinService) GetBinByBinId(c *gin.Context) (*dtos.BinResponse, *common.Error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	bin, err := s.binRepo.GetById(id)
	if err != nil {
		return nil, common.NotFound
	}
	if bin == nil {
		return nil, common.NotFound
	}

	res := modelToBinResponse(bin)
	return &res, nil
}

func (s *BinService) CreateBin(c *gin.Context) (*dtos.BinResponse, *common.Error) {
	var req dtos.BinCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	bin := &model.Bin{
		LocationInWarehouse: req.LocationInWarehouse,
		WarehouseID:         int32(req.WarehouseID),
		StatusInt:           1,
		CreatedAt:           time.Now(),
	}

	err := s.binRepo.Save(bin)
	if err != nil {
		return nil, common.SystemError
	}

	res := modelToBinResponse(bin)
	return &res, nil
}

func (s *BinService) UpdateBin(c *gin.Context) *common.Error {
	var req dtos.BinUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	bin, err := s.binRepo.GetById(req.BinID)
	if err != nil || bin == nil {
		return common.NotFound
	}

	if req.LocationInWarehouse != "" {
		bin.LocationInWarehouse = req.LocationInWarehouse
	}
	if req.WarehouseID != 0 {
		bin.WarehouseID = int32(req.WarehouseID)
	}
	if req.StatusInt != 0 {
		bin.StatusInt = int32(req.StatusInt)
	}
	if req.UpdatedBy != 0 {
		bin.UpdatedBy = int32(req.UpdatedBy)
	}
	bin.UpdatedAt = time.Now()

	err = s.binRepo.Update(bin)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *BinService) DeleteBin(c *gin.Context) *common.Error {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	err := s.binRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func modelToBinResponse(bin *model.Bin) dtos.BinResponse {
	return dtos.BinResponse{
		BinID:               int(bin.BinID),
		LocationInWarehouse: bin.LocationInWarehouse,
		StatusInt:           int(bin.StatusInt),
		WarehouseID:         int (bin.WarehouseID),
		CreatedBy:           int(bin.CreatedBy),
		CreatedAt:           bin.CreatedAt,
		UpdatedBy:           int(bin.UpdatedBy),
		UpdatedAt:           bin.UpdatedAt,
	}
}
