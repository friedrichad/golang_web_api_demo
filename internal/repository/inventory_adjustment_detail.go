package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"gorm.io/gorm"
)

type IInventoryAdjustmentDetail interface {
	IBaseRepository[model.InventoryAdjustmentDetail, int]
	GetByAdjustmentDetailId(adjustmentDetailId int) (*model.InventoryAdjustmentDetail, error)
	GetAllByCondition(query dtos.InventoryAdjustmentDetailRequest) ([]model.InventoryAdjustmentDetail, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryAdjustmentDetail) error
	Update(request *model.InventoryAdjustmentDetail) error
}
type InventoryAdjustmentDetailRepository struct{
	BaseRepository[model.InventoryAdjustmentDetail, int]
	DB *gorm.DB
}
var inventoryAdjustmentDetailRepository IInventoryAdjustmentDetail

func NewInventoryAdjustmentDetailRepository() IInventoryAdjustmentDetail{
	if inventoryAdjustmentDetailRepository == nil{
		inventoryAdjustmentDetailRepository = &InventoryAdjustmentDetailRepository{DB: db.Instance}
		inventoryAdjustmentDetailRepository.SetInstance(db.Instance)	
	}
	return inventoryAdjustmentDetailRepository
}
func (r *InventoryAdjustmentDetailRepository) GetByAdjustmentDetailId(adjustmentDetailId int) (*model.InventoryAdjustmentDetail, error) {
	var inventoryAdjustmentDetail *model.InventoryAdjustmentDetail
	err := r.DB.Where("adjustment_detail_id = ?", adjustmentDetailId).First(&inventoryAdjustmentDetail).Error
	if err != nil {
		return nil, err
	}
	return inventoryAdjustmentDetail, nil
}

func (r *InventoryAdjustmentDetailRepository) GetAllByCondition(query dtos.InventoryAdjustmentDetailRequest) ([]model.InventoryAdjustmentDetail, int, error){
	return r.GetPage("Select iad.* from inventory_adjustment_detail as iad "+
	"where (? is Null or iad.adjustment_detail_id = ?))"+
	"and (? is Null or iad.adjustment_id = ?)) "+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.AdjustmentDetailID, query.AdjustmentDetailID, query.AdjustmentID, query.AdjustmentID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (r *InventoryAdjustmentDetailRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from inventory_adjustment_detail where adjustment_detail_id in ?", ids).Error
}

func (r *InventoryAdjustmentDetailRepository) Save(request *model.InventoryAdjustmentDetail) error {
	return r.BaseRepository.Create(request)
}

func (r *InventoryAdjustmentDetailRepository) Update(request *model.InventoryAdjustmentDetail) error{
	return r.BaseRepository.Update(request)
}
