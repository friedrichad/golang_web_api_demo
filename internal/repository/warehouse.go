package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IWarehouseRepository interface {
	IBaseRepository[model.Warehouse, int]
	GetByWarehouseId(warehouseId int) (*model.Warehouse, error)
	GetAllByCondition(query dtos.WarehouseFilter) ([]model.Warehouse, int, error)
	Delete(ids []int) error
	Save(warehouse *model.Warehouse) error
	Update(warehouse *model.Warehouse) error
	CreateWarehouseTx(warehouse *model.Warehouse) (*model.Warehouse, error)
	UpdateWarehouseTx(warehouse *model.Warehouse) error
	DeleteWarehouseTx(ids []int) error
}

type WarehouseRepository struct {
	BaseRepository[model.Warehouse, int]
	DB *gorm.DB
}

var warehouseRepository IWarehouseRepository

func NewWarehouseRepository() IWarehouseRepository {
	if warehouseRepository == nil {
		warehouseRepository = &WarehouseRepository{DB: db.Instance}
		warehouseRepository.SetInstance(db.Instance)
	}
	return warehouseRepository
}

func (w *WarehouseRepository) GetByWarehouseId(warehouseId int) (*model.Warehouse, error) {
	var warehouse *model.Warehouse
	err := w.DB.Where("warehouse_id = ?", warehouseId).First(&warehouse).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return warehouse, err
}

func (w *WarehouseRepository) GetAllByCondition(query dtos.WarehouseFilter) ([]model.Warehouse, int, error) {
	return w.GetPage("select w.* from warehouse w"+
		" where (? is null or w.warehouse_name like ?)"+
		" and (? is null or w.physical_location = ?)"+
		" and (? is null or w.created_at >= ?)"+
		" and (? is null or w.created_at <= ?)", query.Page, query.Size, query.WarehouseName, query.WarehouseName, query.PhysicalLocation, query.PhysicalLocation, query.DateFrom, query.DateFrom, query.DateTo, query.DateTo)
}

func (w *WarehouseRepository) Delete(ids []int) error {
	return w.DB.Exec("DELETE FROM warehouse WHERE warehouse_id IN ?", ids).Error
}

func (w *WarehouseRepository) Save(warehouse *model.Warehouse) error {
	return w.BaseRepository.Create(warehouse)
}

func (w *WarehouseRepository) Update(warehouse *model.Warehouse) error {
	return w.BaseRepository.Update(warehouse)
}

func (w *WarehouseRepository) WithTx(tx *gorm.DB) *WarehouseRepository {
	return &WarehouseRepository{
		BaseRepository: BaseRepository[model.Warehouse, int]{Instance: tx},
		DB:             tx,
	}
}

// CreateWarehouseTx handles transaction for warehouse creation
func (w *WarehouseRepository) CreateWarehouseTx(warehouse *model.Warehouse) (*model.Warehouse, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	warehouseRepoTx := w.WithTx(tx)
	err := warehouseRepoTx.Save(warehouse)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return warehouse, nil
}

// UpdateWarehouseTx handles transaction for warehouse update
func (w *WarehouseRepository) UpdateWarehouseTx(warehouse *model.Warehouse) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	warehouseRepoTx := w.WithTx(tx)
	err := warehouseRepoTx.Update(warehouse)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteWarehouseTx handles transaction for warehouse deletion
func (w *WarehouseRepository) DeleteWarehouseTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	warehouseRepoTx := w.WithTx(tx)
	err := warehouseRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
