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

func (c *WarehouseController) GetAllWarehouses() gin.HandlerFunc {
	return c.ResponsePage(c.warehouseService.GetAllWarehouses)
}

func (c *WarehouseController) GetWarehouseById() gin.HandlerFunc {
	return c.ResponsePointer(c.warehouseService.GetWarehouseById)
}

func (c *WarehouseController) CreateWarehouse() gin.HandlerFunc {
	return c.ResponsePointer(c.warehouseService.CreateWarehouse)
}

func (c *WarehouseController) UpdateWarehouse() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.warehouseService.UpdateWarehouse)
}

func (c *WarehouseController) DeleteWarehouse() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.warehouseService.DeleteWarehouse)
}
