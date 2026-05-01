package service

import (
	"time"
	"strconv"
	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IWarehouseService interface {
	GetAllWarehouses(c *gin.Context) ([]dtos.WarehouseResponse, int, *common.Error)
	GetWarehouseById(c *gin.Context) (*dtos.WarehouseResponse, *common.Error)
	CreateWarehouse(c *gin.Context) (*dtos.WarehouseResponse, *common.Error)
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

func (s *WarehouseService) GetAllWarehouses(c *gin.Context) ([]dtos.WarehouseResponse, int, *common.Error) {
	var query dtos.WarehouseFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}

	warehouses, total, err := s.warehouseRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	res := make([]dtos.WarehouseResponse, len(warehouses))
	for i, w := range warehouses {
		res[i] = modelToWarehouseResponse(&w)
	}

	return res, total, nil
}

func (s *WarehouseService) GetWarehouseById(c *gin.Context) (*dtos.WarehouseResponse, *common.Error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	warehouse, err := s.warehouseRepo.GetById(id)
	if err != nil || warehouse == nil {
		return nil, common.NotFound
	}

	res := modelToWarehouseResponse(warehouse)
	return &res, nil
}

func (s *WarehouseService) CreateWarehouse(c *gin.Context) (*dtos.WarehouseResponse, *common.Error) {
	var req dtos.WarehouseCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	warehouse := &model.Warehouse{
		WarehouseName:    req.WarehouseName,
		Description:      req.Description,
		PhysicalLocation: req.PhysicalLocation,
		ImageURL:         req.ImageURL,
		CreatedAt:        time.Now(),
	}

	err := s.warehouseRepo.Save(warehouse)
	if err != nil {
		return nil, common.SystemError
	}

	res := modelToWarehouseResponse(warehouse)
	return &res, nil
}

func (s *WarehouseService) UpdateWarehouse(c *gin.Context) *common.Error {
	var req dtos.WarehouseUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	warehouse, err := s.warehouseRepo.GetById(req.WarehouseID)
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
		warehouse.UpdatedBy = int32(req.UpdatedBy)
	}
	warehouse.UpdatedAt = time.Now()

	err = s.warehouseRepo.Update(warehouse)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *WarehouseService) DeleteWarehouse(c *gin.Context) *common.Error {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	err := s.warehouseRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func modelToWarehouseResponse(w *model.Warehouse) dtos.WarehouseResponse {
	return dtos.WarehouseResponse{
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