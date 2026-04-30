package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"gorm.io/gorm"
)

type IInventoryLedger interface {
	IBaseRepository[model.InventoryLedger, int]
	GetByLedgerId(ledgerId int) (*model.InventoryLedger, error)
	GetAllByCondition(query dtos.InventoryLedgerRequest) ([]model.InventoryLedger, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryLedger) error
	Update(request *model.InventoryLedger) error
}

type InventoryLedgerRepository struct{
	BaseRepository[model.InventoryLedger, int]
	DB *gorm.DB
}

var inventoryLedgerRepository IInventoryLedger

func NewInventoryLedgerRepository() IInventoryLedger{
	if inventoryLedgerRepository == nil{
		inventoryLedgerRepository = &InventoryLedgerRepository{DB: db.Instance}
		inventoryLedgerRepository.SetInstance(db.Instance)	
	}
	return inventoryLedgerRepository
}

func (r *InventoryLedgerRepository) GetByLedgerId(ledgerId int) (*model.InventoryLedger, error) {
	var inventoryLedger *model.InventoryLedger
	err := r.DB.Where("ledger_id = ?", ledgerId).First(&inventoryLedger).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return inventoryLedger, err
}

func (r *InventoryLedgerRepository) GetAllByCondition(query dtos.InventoryLedgerRequest) ([]model.InventoryLedger, int, error) {
	return r.GetPage("Select il.* from inventory_ledger as il "+
	"where (? is Null or il.ledger_id = ?))"+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.LedgerID, query.LedgerID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (r *InventoryLedgerRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from inventory_ledger where ledger_id in ?", ids).Error
}

func (r *InventoryLedgerRepository) Save(request *model.InventoryLedger) error {
	return r.BaseRepository.Create(request)
}

func (r *InventoryLedgerRepository) Update(request *model.InventoryLedger) error {
	return r.BaseRepository.Update(request)
}