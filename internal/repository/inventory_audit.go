package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IInventoryAudit interface {
	IBaseRepository[model.InventoryAudit, int]
	GetByRequestId(auditId int) (*model.InventoryAudit, error)
	GetAllByCondition(query dtos.InventoryAuditFilter) ([]model.InventoryAudit, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryAudit) error
	Update(request *model.InventoryAudit) error
	CreateInventoryAuditTx(audit *model.InventoryAudit, details []model.InventoryAuditDetail) (*model.InventoryAudit, error)
	CreateInventoryAuditDetailsTx(details []model.InventoryAuditDetail) error
	UpdateInventoryAuditTx(audit *model.InventoryAudit) error
	DeleteInventoryAuditTx(ids []int) error
}

type InventoryAuditRepository struct {
	BaseRepository[model.InventoryAudit, int]
	DB *gorm.DB
}

var inventoryAuditRepository IInventoryAudit

func NewInventoryAuditRepository() IInventoryAudit {
	if inventoryAuditRepository == nil {
		inventoryAuditRepository = &InventoryAuditRepository{DB: db.Instance}
		inventoryAuditRepository.SetInstance(db.Instance)
	}
	return inventoryAuditRepository
}

func (r *InventoryAuditRepository) GetByRequestId(auditId int) (*model.InventoryAudit, error) {
	var inventoryAudit *model.InventoryAudit
	err := r.DB.Where("audit_id = ?", auditId).First(&inventoryAudit).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return inventoryAudit, err
}
func (r *InventoryAuditRepository) GetAllByCondition(query dtos.InventoryAuditFilter) ([]model.InventoryAudit, int, error) {
	return r.GetPage("Select ia.* from inventory_audit as ia "+
		" where (? is Null or ia.audit_id = ?)"+
		" and (? is null or ia.warehouse_id = ?) "+
		" and (? is null or ia.created_at >= ?) "+
		" and (? is null or ia.created_at < ?) ", query.Page, query.Size, query.AuditID, query.AuditID, query.WarehouseID, query.WarehouseID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (r *InventoryAuditRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from inventory_audit where audit_id in ?", ids).Error
}
func (r *InventoryAuditRepository) Save(request *model.InventoryAudit) error {
	return r.BaseRepository.Create(request)
}
func (r *InventoryAuditRepository) Update(request *model.InventoryAudit) error {
	return r.BaseRepository.Update(request)
}

func (r *InventoryAuditRepository) WithTx(tx *gorm.DB) *InventoryAuditRepository {
	return &InventoryAuditRepository{
		BaseRepository: BaseRepository[model.InventoryAudit, int]{Instance: tx},
		DB:             tx,
	}
}

// CreateInventoryAuditTx handles transaction for inventory audit creation with details
func (r *InventoryAuditRepository) CreateInventoryAuditTx(audit *model.InventoryAudit, details []model.InventoryAuditDetail) (*model.InventoryAudit, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	auditRepoTx := r.WithTx(tx)
	detailRepoTx := NewInventoryAuditDetailRepository().(*InventoryAuditDetailRepository).WithTx(tx)

	err := auditRepoTx.Save(audit)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, detail := range details {
		detail.AuditID = audit.AuditID
		err := detailRepoTx.Save(&detail)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return audit, nil
}

// CreateInventoryAuditDetailsTx handles transaction for creating multiple audit details
func (r *InventoryAuditRepository) CreateInventoryAuditDetailsTx(details []model.InventoryAuditDetail) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	detailRepoTx := NewInventoryAuditDetailRepository().(*InventoryAuditDetailRepository).WithTx(tx)

	for _, detail := range details {
		err := detailRepoTx.Save(&detail)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// UpdateInventoryAuditTx handles transaction for inventory audit update
func (r *InventoryAuditRepository) UpdateInventoryAuditTx(audit *model.InventoryAudit) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	auditRepoTx := r.WithTx(tx)
	err := auditRepoTx.Update(audit)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteInventoryAuditTx handles transaction for inventory audit deletion
func (r *InventoryAuditRepository) DeleteInventoryAuditTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	auditRepoTx := r.WithTx(tx)
	err := auditRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
