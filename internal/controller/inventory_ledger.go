package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IInventoryLedgerController interface {
	GetAllLedgers() gin.HandlerFunc
	GetLedgerById() gin.HandlerFunc
	ExportLedgersExcel() gin.HandlerFunc
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

func (c *InventoryLedgerController) ExportLedgersExcel() gin.HandlerFunc {
	return func(g *gin.Context) {
		fileBytes, err := c.ledgerService.ExportInventoryLedgersToExcel(g)
		if err != nil {
			c.Error(g, err, nil)
			return
		}

		fileName := fmt.Sprintf("inventory_ledger_%s.xlsx", time.Now().Format("20060102_150405"))
		g.Header("Content-Description", "File Transfer")
		g.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		g.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
	}
}
