package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IInventoryLedgerController interface {
	GetAllLedgers() gin.HandlerFunc
	GetLedgerById() gin.HandlerFunc
	ExportLedgersExcel() gin.HandlerFunc
}

type InventoryLedgerController struct {
	BaseController[model.InventoryLedgerResponse]
	ledgerService service.IInventoryLedgerService
}

func NewInventoryLedgerController() IInventoryLedgerController {
	return &InventoryLedgerController{
		ledgerService: service.NewInventoryLedgerService(),
	}
}

// GetAllLedgers godoc
// @Summary Get all inventory ledgers
// @Description Get a list of all inventory ledgers with optional filtering and pagination
// @Tags ledgers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.InventoryLedgerResponse
// @Failure 400 {object} common.Error
// @Router /ledgers [get]
func (c *InventoryLedgerController) GetAllLedgers() gin.HandlerFunc {
	return c.ResponsePage(c.ledgerService.GetAllInventoryLedgers)
}

// GetLedgerById godoc
// @Summary Get a ledger by ID
// @Description Retrieve details of a specific inventory ledger by its ID
// @Tags ledgers
// @Accept json
// @Produce json
// @Param id path int true "Ledger ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.InventoryLedgerResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /ledgers/{id} [get]
func (c *InventoryLedgerController) GetLedgerById() gin.HandlerFunc {
	return c.ResponsePointer(c.ledgerService.GetInventoryLedgerById)
}

// ExportLedgersExcel godoc
// @Summary Export inventory ledgers to Excel
// @Description Export all inventory ledgers to an Excel file
// @Tags ledgers
// @Accept json
// @Produce octet-stream
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} common.Error
// @Router /ledgers/export [get]
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
