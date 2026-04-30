package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"gorm.io/gorm"
)

type IInventoryAudit interface {
	IBaseRepository[model.InventoryAudit, int]
	GetByRequestId(auditId int) (*model.InventoryAudit, error)
	GetAllByCondition(query dtos.InventoryAuditRequest) ([]model.InventoryAudit, int, error)
	Delete(ids []int) error
	Save(request *model.InventoryAudit) error
	Update(request *model.InventoryAudit) error
}

type InventoryAuditRepository struct{
	BaseRepository[model.InventoryAudit, int]
	DB *gorm.DB
}
var inventoryAuditRepository IInventoryAudit

func NewInventoryAuditRepository() IInventoryAudit{
	if inventoryAuditRepository == nil{
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
func (r *InventoryAuditRepository) GetAllByCondition(query dtos.InventoryAuditRequest) ([]model.InventoryAudit, int, error) {
	return r.GetPage("Select ia.* from inventory_audit as ia "+
	"where (? is Null or ia.audit_id = ?))"+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.AuditID, query.AuditID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
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