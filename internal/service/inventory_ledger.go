package service

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type IInventoryLedgerService interface {
	GetAllInventoryLedgers(c *gin.Context) ([]dtos.InventoryLedgerResponse, int, *common.Error)
	GetInventoryLedgerById(c *gin.Context) (*dtos.InventoryLedgerResponse, *common.Error)
	ExportInventoryLedgersToExcel(c *gin.Context) ([]byte, *common.Error)
	// Internal method for creating ledger entries from other services
	CreateInventoryLedgerEntry(req *dtos.InventoryLedgerCreate) error
}

type InventoryLedgerService struct {
	ledgerRepo repository.IInventoryLedger
}

var inventoryLedgerService IInventoryLedgerService

func NewInventoryLedgerService() IInventoryLedgerService {
	if inventoryLedgerService == nil {
		inventoryLedgerService = &InventoryLedgerService{
			ledgerRepo: repository.NewInventoryLedgerRepository(),
		}
	}
	return inventoryLedgerService
}

func (s *InventoryLedgerService) GetAllInventoryLedgers(c *gin.Context) ([]dtos.InventoryLedgerResponse, int, *common.Error) {
	var query dtos.InventoryLedgerFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	ledgers, total, err := s.ledgerRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	ledgerResponses := make([]dtos.InventoryLedgerResponse, len(ledgers))
	for i, ledger := range ledgers {
		ledgerResponses[i] = modelToInventoryLedgerResponse(&ledger)
	}

	return ledgerResponses, total, nil
}

func (s *InventoryLedgerService) GetInventoryLedgerById(c *gin.Context) (*dtos.InventoryLedgerResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	ledgerId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	ledger, err := s.ledgerRepo.GetByLedgerId(int(ledgerId))
	if err != nil {
		return nil, common.NotFound
	}

	if ledger == nil {
		return nil, &common.Error{Code: "404", Message: "Bản ghi không tồn tại"}
	}

	ledgerResponse := modelToInventoryLedgerResponse(ledger)
	return &ledgerResponse, nil
}

// ExportInventoryLedgersToExcel - Export inventory ledger data to Excel file
func (s *InventoryLedgerService) ExportInventoryLedgersToExcel(c *gin.Context) ([]byte, *common.Error) {
	var query dtos.InventoryLedgerFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, common.RequestInvalid
	}

	ledgers, total, err := s.ledgerRepo.GetAllByCondition(query)
	if err != nil {
		return nil, common.SystemError
	}
	if total == 0 {
		return nil, common.NotFound
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Inventory Ledger"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// --- Header style ---
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#D9E2F3", Style: 1},
			{Type: "top", Color: "#D9E2F3", Style: 1},
			{Type: "bottom", Color: "#D9E2F3", Style: 1},
			{Type: "right", Color: "#D9E2F3", Style: 1},
		},
	})

	// --- Data style ---
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#D9E2F3", Style: 1},
			{Type: "top", Color: "#D9E2F3", Style: 1},
			{Type: "bottom", Color: "#D9E2F3", Style: 1},
			{Type: "right", Color: "#D9E2F3", Style: 1},
		},
	})

	// --- Number style (for quantity columns) ---
	numberStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		NumFmt: 4, // #,##0.00
		Border: []excelize.Border{
			{Type: "left", Color: "#D9E2F3", Style: 1},
			{Type: "top", Color: "#D9E2F3", Style: 1},
			{Type: "bottom", Color: "#D9E2F3", Style: 1},
			{Type: "right", Color: "#D9E2F3", Style: 1},
		},
	})

	// --- Date style ---
	dateStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		NumFmt: 22, // m/d/yy h:mm
		Border: []excelize.Border{
			{Type: "left", Color: "#D9E2F3", Style: 1},
			{Type: "top", Color: "#D9E2F3", Style: 1},
			{Type: "bottom", Color: "#D9E2F3", Style: 1},
			{Type: "right", Color: "#D9E2F3", Style: 1},
		},
	})

	// Define headers
	headers := []string{
		"STT", "Ledger ID", "Component ID", "Warehouse ID", "Bin ID",
		"Loại tham chiếu", "Mã tham chiếu", "Mô tả",
		"SL thay đổi", "SL sau", "Ghi chú", "Người tạo", "Ngày tạo",
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 6, "B": 12, "C": 14, "D": 14, "E": 10,
		"F": 18, "G": 14, "H": 30,
		"I": 14, "J": 14, "K": 25, "L": 12, "M": 20,
	}
	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write data rows
	for rowIdx, ledger := range ledgers {
		row := rowIdx + 2 // data starts from row 2
		refTypeName := constants.GetLedgerReferenceTypeName(ledger.ReferenceType)

		// STT
		cellA, _ := excelize.CoordinatesToCellName(1, row)
		f.SetCellValue(sheetName, cellA, rowIdx+1)
		f.SetCellStyle(sheetName, cellA, cellA, dataStyle)

		// Ledger ID
		cellB, _ := excelize.CoordinatesToCellName(2, row)
		f.SetCellValue(sheetName, cellB, ledger.LedgerID)
		f.SetCellStyle(sheetName, cellB, cellB, dataStyle)

		// Component ID
		cellC, _ := excelize.CoordinatesToCellName(3, row)
		f.SetCellValue(sheetName, cellC, ledger.ComponentID)
		f.SetCellStyle(sheetName, cellC, cellC, dataStyle)

		// Warehouse ID
		cellD, _ := excelize.CoordinatesToCellName(4, row)
		f.SetCellValue(sheetName, cellD, ledger.WarehouseID)
		f.SetCellStyle(sheetName, cellD, cellD, dataStyle)

		// Bin ID
		cellE, _ := excelize.CoordinatesToCellName(5, row)
		f.SetCellValue(sheetName, cellE, ledger.BinID)
		f.SetCellStyle(sheetName, cellE, cellE, dataStyle)

		// Loại tham chiếu
		cellF, _ := excelize.CoordinatesToCellName(6, row)
		f.SetCellValue(sheetName, cellF, refTypeName)
		f.SetCellStyle(sheetName, cellF, cellF, dataStyle)

		// Mã tham chiếu
		cellG, _ := excelize.CoordinatesToCellName(7, row)
		f.SetCellValue(sheetName, cellG, ledger.ReferenceTypeID)
		f.SetCellStyle(sheetName, cellG, cellG, dataStyle)

		// Mô tả
		cellH, _ := excelize.CoordinatesToCellName(8, row)
		f.SetCellValue(sheetName, cellH, ledger.Description)
		f.SetCellStyle(sheetName, cellH, cellH, dataStyle)

		// SL thay đổi
		cellI, _ := excelize.CoordinatesToCellName(9, row)
		f.SetCellValue(sheetName, cellI, ledger.QuantityChange)
		f.SetCellStyle(sheetName, cellI, cellI, numberStyle)

		// SL sau
		cellJ, _ := excelize.CoordinatesToCellName(10, row)
		f.SetCellValue(sheetName, cellJ, ledger.QuantityAfter)
		f.SetCellStyle(sheetName, cellJ, cellJ, numberStyle)

		// Ghi chú
		cellK, _ := excelize.CoordinatesToCellName(11, row)
		f.SetCellValue(sheetName, cellK, ledger.Note)
		f.SetCellStyle(sheetName, cellK, cellK, dataStyle)

		// Người tạo
		cellL, _ := excelize.CoordinatesToCellName(12, row)
		f.SetCellValue(sheetName, cellL, ledger.CreatedBy)
		f.SetCellStyle(sheetName, cellL, cellL, dataStyle)

		// Ngày tạo
		cellM, _ := excelize.CoordinatesToCellName(13, row)
		f.SetCellValue(sheetName, cellM, ledger.CreatedAt.Format("02/01/2006 15:04:05"))
		f.SetCellStyle(sheetName, cellM, cellM, dateStyle)
	}

	// Set row height for header
	f.SetRowHeight(sheetName, 1, 25)

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, common.SystemError
	}

	return buf.Bytes(), nil
}

func modelToInventoryLedgerResponse(ledger *model.InventoryLedger) dtos.InventoryLedgerResponse {
	return dtos.InventoryLedgerResponse{
		LedgerID:        int(ledger.LedgerID),
		ComponentID:     int(ledger.ComponentID),
		WarehouseID:     int(ledger.WarehouseID),
		BinID:           int(ledger.BinID),
		ReferenceType:   int(ledger.ReferenceType),
		ReferenceTypeID: int(ledger.ReferenceTypeID),
		Description:     ledger.Description,
		QuantityChange:  ledger.QuantityChange,
		QuantityAfter:   ledger.QuantityAfter,
		Note:            ledger.Note,
		CreatedAt:       ledger.CreatedAt,
		CreatedBy:       int(ledger.CreatedBy),
	}
}

// CreateInventoryLedgerEntry - Internal method to create ledger entries from other services
func (s *InventoryLedgerService) CreateInventoryLedgerEntry(req *dtos.InventoryLedgerCreate) error {
	// Validate reference type
	if !constants.IsValidLedgerReferenceType(req.ReferenceType) {
		return fmt.Errorf("loại tham chiếu không hợp lệ")
	}

	// Create ledger model
	ledger := &model.InventoryLedger{
		ComponentID:     req.ComponentID,
		WarehouseID:     req.WarehouseID,
		BinID:           req.BinID,
		ReferenceType:   req.ReferenceType,
		ReferenceTypeID: req.ReferenceTypeID,
		Description:     req.Description,
		QuantityChange:  req.QuantityChange,
		QuantityAfter:   req.QuantityAfter,
		Note:            req.Note,
		CreatedBy:       req.CreatedBy,
		CreatedAt:       time.Now(),
		UpdatedBy:       req.CreatedBy,
		UpdatedAt:       time.Now(),
	}

	// Save ledger entry
	if err := s.ledgerRepo.Save(ledger); err != nil {
		return err
	}

	return nil
}
