package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IInventoryLedgerController interface {
	GetAllLedgers() gin.HandlerFunc
	GetLedgerById() gin.HandlerFunc
}

type InventoryLedgerController struct {
	BaseController[dtos.InventoryLedgerResponse]
	ledgerService service.IInventoryLedgerService
}

func NewInventoryLedgerController() IInventoryLedgerController {
	return &InventoryLedgerController{
		ledgerService: service.NewInventoryLedgerService(),
	}
}

func (c *InventoryLedgerController) GetAllLedgers() gin.HandlerFunc {
	return c.ResponsePage(c.ledgerService.GetAllInventoryLedgers)
}

func (c *InventoryLedgerController) GetLedgerById() gin.HandlerFunc {
	return c.ResponsePointer(c.ledgerService.GetInventoryLedgerById)
}
