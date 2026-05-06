package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"gorm.io/gorm"
)

type IInventoryAdjustment interface {
	IBaseRepository[model.InventoryAdjustment, int]
	GetByAuditId(auditId int) (*model.InventoryAdjustment, error)
	GetByAdjustmentId(adjustmentId int) (*model.InventoryAdjustment, error)
	GetAllByCondition(query dtos.InventoryAdjustmentFilter) ([]model.InventoryAdjustment, int, error)
	Delete(ids []int) error
	DeleteIfPending(ids []int) error
	Save(request *model.InventoryAdjustment) error
	Update(request *model.InventoryAdjustment) error
	CreateWithDetails(adjustment *model.InventoryAdjustment, details []model.InventoryAdjustmentDetail) error
	ApproveWithComponentBinUpdate(adjustment *model.InventoryAdjustment, details []model.InventoryAdjustmentDetail, updatedBy int) error
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

func (r *InventoryAdjustmentRepository) CreateWithDetails(adjustment *model.InventoryAdjustment, details []model.InventoryAdjustmentDetail) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Save adjustment
	if err := tx.Create(adjustment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Save details if any
	if len(details) > 0 {
		for i := range details {
			details[i].AdjustmentID = adjustment.AdjustmentID
			if err := tx.Create(&details[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *InventoryAdjustmentRepository) ApproveWithComponentBinUpdate(
	adjustment *model.InventoryAdjustment,
	details []model.InventoryAdjustmentDetail,
	updatedBy int) error {

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Update adjustment status
	adjustment.StatusInt = constants.InventoryAdjustmentStatusApproved
	adjustment.ApprovedID = updatedBy
	adjustment.ApprovedTime = adjustment.ApprovedTime
	adjustment.UpdatedAt = adjustment.UpdatedAt
	adjustment.UpdatedBy = updatedBy

	if err := tx.Model(adjustment).Updates(adjustment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update component bins for each detail
	for _, detail := range details {
		compBin := &model.ComponentBin{}
		result := tx.Where("component_id = ? AND bin_id = ?", detail.ComponentID, detail.BinID).First(compBin)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			tx.Rollback()
			return result.Error
		}

		if result.Error == gorm.ErrRecordNotFound {
			// Create if not exists
			compBin = &model.ComponentBin{
				ComponentID: detail.ComponentID,
				BinID:       detail.BinID,
				Quantity:    detail.QuantityAfter,
				CreatedAt:   adjustment.CreatedAt,
				CreatedBy:   updatedBy,
			}
			if err := tx.Create(compBin).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Update existing
			compBin.Quantity = detail.QuantityAfter
			compBin.UpdatedAt = adjustment.UpdatedAt
			compBin.UpdatedBy = updatedBy
			if err := tx.Model(compBin).Updates(compBin).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *InventoryAdjustmentRepository) DeleteIfPending(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Check if all adjustments are in pending status
	for _, id := range ids {
		adjustment := &model.InventoryAdjustment{}
		result := tx.Where("adjustment_id = ?", id).First(adjustment)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			tx.Rollback()
			return result.Error
		}

		if result.Error == nil && adjustment.StatusInt != constants.InventoryAdjustmentStatusPending {
			tx.Rollback()
			return gorm.ErrInvalidData // Use this to indicate validation error
		}
	}

	// All are pending, so delete them
	if err := tx.Exec("delete from inventory_adjustment where adjustment_id in ?", ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
