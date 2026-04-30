package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"gorm.io/gorm"
)

type IInventoryAuditDetail interface {
	IBaseRepository[model.InventoryAuditDetail, int]
	GetByInventoryAuditDetailId(auditDetailId int) (*model.InventoryAuditDetail, error)
	GetAllByCondition(query dtos.InventoryAuditDetailRequest) ([]model.InventoryAuditDetail, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryAuditDetail) error
	Update(request *model.InventoryAuditDetail) error
}

type InventoryAuditDetailRepository struct{
	BaseRepository[model.InventoryAuditDetail, int]
	DB *gorm.DB
}
var inventoryAuditDetailRepository IInventoryAuditDetail

func NewInventoryAuditDetailRepository() IInventoryAuditDetail{
	if inventoryAuditDetailRepository == nil{
		inventoryAuditDetailRepository = &InventoryAuditDetailRepository{DB: db.Instance}
		inventoryAuditDetailRepository.SetInstance(db.Instance)	
	}
	return inventoryAuditDetailRepository
}
func (r *InventoryAuditDetailRepository) GetByInventoryAuditDetailId(auditDetailId int) (*model.InventoryAuditDetail, error) {
	var inventoryAuditDetail *model.InventoryAuditDetail
	err := r.DB.Where("id = ?", auditDetailId).First(&inventoryAuditDetail).Error
	if err != nil {
		return nil, err
	}
	return inventoryAuditDetail, nil
}
func (r *InventoryAuditDetailRepository) GetAllByCondition(query dtos.InventoryAuditDetailRequest) ([]model.InventoryAuditDetail, int, error) {
	return r.GetPage("Select iad.* from inventory_audit_detail as iad "+
	"where (? is Null or iad.audit_detail_id = ?))"+
	"and (? is Null or iad.audit_id = ?)) "+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.AuditDetailID, query.AuditID, query.AuditID, query.AuditID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (r *InventoryAuditDetailRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from inventory_audit_detail where audit_detail_id in ?", ids).Error
}
func (r *InventoryAuditDetailRepository) Save(request *model.InventoryAuditDetail) error {
	return r.BaseRepository.Create(request)
}
func (r *InventoryAuditDetailRepository) Update(request *model.InventoryAuditDetail) error {
	return r.BaseRepository.Update(request)
}