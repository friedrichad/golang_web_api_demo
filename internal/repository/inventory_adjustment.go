package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IInventoryAdjustment interface {
	IBaseRepository[model.InventoryAdjustment, int]
	GetByAuditId(auditId int) (*model.InventoryAdjustment, error)
	GetByAdjustmentId(adjustmentId int) (*model.InventoryAdjustment, error)
	GetAllByCondition(query dtos.InventoryAdjustmentFilter) ([]model.InventoryAdjustment, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryAdjustment) error
	Update(request *model.InventoryAdjustment) error
	CreateInventoryAdjustmentTx(adjustment *model.InventoryAdjustment, details []model.InventoryAdjustmentDetail) (*model.InventoryAdjustment, error)
	UpdateInventoryAdjustmentTx(adjustment *model.InventoryAdjustment) error
	DeleteInventoryAdjustmentTx(ids []int) error
}

type InventoryAdjustmentRepository struct {
	BaseRepository[model.InventoryAdjustment, int]
	DB *gorm.DB
}

var inventoryAdjustmentRepository IInventoryAdjustment

func NewInventoryAdjustmentRepository() IInventoryAdjustment {
	if inventoryAdjustmentRepository == nil {
		inventoryAdjustmentRepository = &InventoryAdjustmentRepository{DB: db.Instance}
		inventoryAdjustmentRepository.SetInstance(db.Instance)
	}
	return inventoryAdjustmentRepository
}
func (r *InventoryAdjustmentRepository) GetByAuditId(auditId int) (*model.InventoryAdjustment, error) {
	var inventoryAdjustment *model.InventoryAdjustment
	err := r.DB.Where("audit_id = ?", auditId).First(&inventoryAdjustment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return inventoryAdjustment, err
}

func (r *InventoryAdjustmentRepository) GetByAdjustmentId(adjustmentId int) (*model.InventoryAdjustment, error) {
	var inventoryAdjustment *model.InventoryAdjustment
	err := r.DB.Where("adjustment_id = ?", adjustmentId).First(&inventoryAdjustment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return inventoryAdjustment, err
}

func (r *InventoryAdjustmentRepository) GetAllByCondition(query dtos.InventoryAdjustmentFilter) ([]model.InventoryAdjustment, int, error) {
	return r.GetPage("Select ia.* from inventory_adjustment as ia "+
		" where (? is Null or ia.adjustment_id = ?)"+
		" and (? is null or ia.created_at >= ?) "+
		" and (? is null or ia.created_at < ?) ", query.Page, query.Size, query.AdjustmentID, query.AdjustmentID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (r *InventoryAdjustmentRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from inventory_adjustment where adjustment_id in ?", ids).Error
}
func (r *InventoryAdjustmentRepository) Save(request *model.InventoryAdjustment) error {
	return r.BaseRepository.Create(request)
}
func (r *InventoryAdjustmentRepository) Update(request *model.InventoryAdjustment) error {
	return r.BaseRepository.Update(request)
}

func (r *InventoryAdjustmentRepository) WithTx(tx *gorm.DB) *InventoryAdjustmentRepository {
	return &InventoryAdjustmentRepository{
		BaseRepository: BaseRepository[model.InventoryAdjustment, int]{Instance: tx},
		DB:             tx,
	}
}

// CreateInventoryAdjustmentTx handles transaction for inventory adjustment creation with details
func (r *InventoryAdjustmentRepository) CreateInventoryAdjustmentTx(adjustment *model.InventoryAdjustment, details []model.InventoryAdjustmentDetail) (*model.InventoryAdjustment, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	adjustmentRepoTx := r.WithTx(tx)
	detailRepoTx := NewInventoryAdjustmentDetailRepository().(*InventoryAdjustmentDetailRepository).WithTx(tx)

	err := adjustmentRepoTx.Save(adjustment)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, detail := range details {
		detail.AdjustmentID = adjustment.AdjustmentID
		err := detailRepoTx.Save(&detail)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return adjustment, nil
}

// UpdateInventoryAdjustmentTx handles transaction for inventory adjustment update
func (r *InventoryAdjustmentRepository) UpdateInventoryAdjustmentTx(adjustment *model.InventoryAdjustment) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	adjustmentRepoTx := r.WithTx(tx)
	err := adjustmentRepoTx.Update(adjustment)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteInventoryAdjustmentTx handles transaction for inventory adjustment deletion
func (r *InventoryAdjustmentRepository) DeleteInventoryAdjustmentTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	adjustmentRepoTx := r.WithTx(tx)
	err := adjustmentRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
