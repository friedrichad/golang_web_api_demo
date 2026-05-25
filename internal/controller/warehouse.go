package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IWarehouseController interface {
	GetAllWarehouses() gin.HandlerFunc
	GetWarehouseById() gin.HandlerFunc
	CreateWarehouse() gin.HandlerFunc
	UpdateWarehouse() gin.HandlerFunc
	DeleteWarehouse() gin.HandlerFunc
}

type WarehouseController struct {
	BaseController[model.WarehouseResponse]
	warehouseService service.IWarehouseService
}

func NewWarehouseController() IWarehouseController {
	warehouseService := service.NewWarehouseService()
	return &WarehouseController{warehouseService: warehouseService}
}

// GetAllWarehouses godoc
// @Summary Get all warehouses
// @Description Get a list of all warehouses with optional filtering and pagination
// @Tags warehouses
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.WarehouseResponse
// @Failure 400 {object} common.Error
// @Router /warehouses [get]
func (c *WarehouseController) GetAllWarehouses() gin.HandlerFunc {
	return c.ResponsePage(c.warehouseService.GetAllWarehouses)
}

// GetWarehouseById godoc
// @Summary Get a warehouse by ID
// @Description Retrieve details of a specific warehouse by its ID
// @Tags warehouses
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.WarehouseResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /warehouses/{id} [get]
func (c *WarehouseController) GetWarehouseById() gin.HandlerFunc {
	return c.ResponsePointer(c.warehouseService.GetWarehouseById)
}

// CreateWarehouse godoc
// @Summary Create a new warehouse
// @Description Add a new warehouse to the system
// @Tags warehouses
// @Accept json
// @Produce json
// @Param request body model.WarehouseCreate true "Warehouse Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.WarehouseResponse
// @Failure 400 {object} common.Error
// @Router /warehouses [post]
func (c *WarehouseController) CreateWarehouse() gin.HandlerFunc {
	return c.ResponsePointer(c.warehouseService.CreateWarehouse)
}

// UpdateWarehouse godoc
// @Summary Update an existing warehouse
// @Description Modify details of an existing warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param request body model.WarehouseUpdate true "Warehouse Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /warehouses [put]
func (c *WarehouseController) UpdateWarehouse() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.warehouseService.UpdateWarehouse)
}

// DeleteWarehouse godoc
// @Summary Delete warehouses
// @Description Delete one or multiple warehouses by their IDs
// @Tags warehouses
// @Accept json
// @Produce json
// @Param request body []int true "Array of Warehouse IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /warehouses [delete]
func (c *WarehouseController) DeleteWarehouse() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.warehouseService.DeleteWarehouse)
}
