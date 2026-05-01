package service

import (
	"strconv"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IInventoryLedgerService interface {
	GetAllInventoryLedgers(c *gin.Context) ([]dtos.InventoryLedgerResponse, int, *common.Error)
	GetInventoryLedgerById(c *gin.Context) (*dtos.InventoryLedgerResponse, *common.Error)
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

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
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

	ledgerId, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return nil, common.RequestInvalid
	}

	ledger, err := s.ledgerRepo.GetById(int32(ledgerId))
	if err != nil {
		return nil, common.NotFound
	}

	if ledger == nil {
		return nil, &common.Error{Code: "404", Message: "Bản ghi không tồn tại"}
	}

	ledgerResponse := modelToInventoryLedgerResponse(ledger)
	return &ledgerResponse, nil
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
