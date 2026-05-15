package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IWarehouseService interface {
	GetAllWarehouses(c *gin.Context) ([]model.WarehouseResponse, int, *common.Error)
	GetWarehouseById(c *gin.Context) (*model.WarehouseResponse, *common.Error)
	CreateWarehouse(c *gin.Context) (*model.WarehouseResponse, *common.Error)
	UpdateWarehouse(c *gin.Context) *common.Error
	DeleteWarehouse(c *gin.Context) *common.Error
}

type WarehouseService struct {
	warehouseRepo repository.IWarehouseRepository
}

var warehouseService IWarehouseService

func NewWarehouseService() IWarehouseService {
	if warehouseService == nil {
		warehouseService = &WarehouseService{
			warehouseRepo: repository.NewWarehouseRepository(),
		}
	}
	return warehouseService
}

func (s *WarehouseService) GetAllWarehouses(c *gin.Context) ([]model.WarehouseResponse, int, *common.Error) {
	var query model.WarehouseFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	warehouses, total, err := s.warehouseRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	res := make([]model.WarehouseResponse, len(warehouses))
	for i, w := range warehouses {
		res[i] = modelToWarehouseResponse(&w)
	}

	return res, total, nil
}

func (s *WarehouseService) GetWarehouseById(c *gin.Context) (*model.WarehouseResponse, *common.Error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	warehouse, err := s.warehouseRepo.GetByWarehouseId(id)
	if err != nil || warehouse == nil {
		return nil, common.NotFound
	}

	res := modelToWarehouseResponse(warehouse)
	return &res, nil
}

func (s *WarehouseService) CreateWarehouse(c *gin.Context) (*model.WarehouseResponse, *common.Error) {
	var req model.WarehouseCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	warehouse := &model.Warehouse{
		WarehouseName:    req.WarehouseName,
		Description:      req.Description,
		PhysicalLocation: req.PhysicalLocation,
		ImageURL:         req.ImageURL,
		CreatedAt:        time.Now(),
	}

	result, err := s.warehouseRepo.CreateWarehouseTx(warehouse)
	if err != nil {
		return nil, common.SystemError
	}

	res := modelToWarehouseResponse(result)
	return &res, nil
}

func (s *WarehouseService) UpdateWarehouse(c *gin.Context) *common.Error {
	var req model.WarehouseUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return &common.Error{Code: "400", Message: err.Error()}
	}

	warehouse, err := s.warehouseRepo.GetByWarehouseId(req.WarehouseID)
	if err != nil || warehouse == nil {
		return common.NotFound
	}

	if req.WarehouseName != "" {
		warehouse.WarehouseName = req.WarehouseName
	}
	if req.Description != "" {
		warehouse.Description = req.Description
	}
	if req.PhysicalLocation != "" {
		warehouse.PhysicalLocation = req.PhysicalLocation
	}
	if req.ImageURL != "" {
		warehouse.ImageURL = req.ImageURL
	}
	if req.UpdatedBy != 0 {
		warehouse.UpdatedBy = int(req.UpdatedBy)
	}
	warehouse.UpdatedAt = time.Now()

	if err := s.warehouseRepo.UpdateWarehouseTx(warehouse); err != nil {
		return common.SystemError
	}

	return nil
}

func (s *WarehouseService) DeleteWarehouse(c *gin.Context) *common.Error {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	if err := s.warehouseRepo.DeleteWarehouseTx(ids); err != nil {
		return common.SystemError
	}

	return nil
}

func modelToWarehouseResponse(w *model.Warehouse) model.WarehouseResponse {
	return model.WarehouseResponse{
		WarehouseID:      int(w.WarehouseID),
		WarehouseName:    w.WarehouseName,
		Description:      w.Description,
		PhysicalLocation: w.PhysicalLocation,
		ImageURL:         w.ImageURL,
		CreatedBy:        int(w.CreatedBy),
		CreatedAt:        w.CreatedAt,
		UpdatedBy:        int(w.UpdatedBy),
		UpdatedAt:        w.UpdatedAt,
	}
}
