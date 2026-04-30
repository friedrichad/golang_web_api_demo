package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IInventoryAdjustment interface {
	IBaseRepository[model.InventoryAdjustment, int]
	GetByAdjustmentId(adjustmentId int) (*model.InventoryAdjustment, error)
	GetAllByCondition(query dtos.InventoryAdjustmentFilter) ([]model.InventoryAdjustment, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryAdjustment) error
	Update(request *model.InventoryAdjustment) error
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
		"where (? is Null or ia.adjustment_id = ?))"+
		"and (? is null or create_at >= ?) "+
		"and (? is null or create_at < ?) ", query.Page, query.Size, query.AdjustmentID, query.AdjustmentID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
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
