package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"gorm.io/gorm"
)

type IRequestPermissionRepository interface {
	IBaseRepository[model.RequestPermission, int]
	GetRequestPermissionByRequestId(requestId int) ([]model.RequestPermission, error)
	GetByRequestPermissionId(requestPermissionId int) (*model.RequestPermission, error)
	GetAllByCondition(query dtos.RequestPermissionFilter) ([]model.RequestPermission, int, error)
	Delete(ids []int) error
	Save(requestPermission *model.RequestPermission) error
	Update(requestPermission *model.RequestPermission) error
}

type RequestPermissionRepository struct {
	BaseRepository[model.RequestPermission, int]
	DB *gorm.DB
}

var requestPermissionRepository IRequestPermissionRepository

func NewRequestPermissionRepository() IRequestPermissionRepository {
	if requestPermissionRepository == nil {
		requestPermissionRepository = &RequestPermissionRepository{DB: db.Instance}
		requestPermissionRepository.SetInstance(db.Instance)
	}
	return requestPermissionRepository
}

func (r * RequestPermissionRepository) GetByRequestPermissionId(requestPermissionId int) (*model.RequestPermission, error) {
	var requestPermission *model.RequestPermission
	err := r.DB.Where("request_permission_id = ?", requestPermissionId).First(&requestPermission).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return requestPermission, err
}

func (r *RequestPermissionRepository) GetRequestPermissionByRequestId(requestId int) ([]model.RequestPermission, error) {
	var requestPermissions []model.RequestPermission
	err := r.DB.Where("request_id = ?", requestId).Find(&requestPermissions).Error
	return requestPermissions, err
}

func (r *RequestPermissionRepository) GetAllByCondition(query dtos.RequestPermissionFilter) ([]model.RequestPermission, int, error){
	return r.GetPage("Select * from RequestPermission" +
		" where (? is null or request_permission_id = ?)"+
		" and (? is null or request_id = ?)" +
		" and (? is null or menu_id =?)"+
		" and (? is null or permission_id = ?)" +
		" and (? is null or created_at >= ?)"+
		" and (? is null or created_at < ?)", query.Page, query.Size, query.RequestPermissionID, query.RequestPermissionID, query.RequestID, query.RequestID, query.MenuID, query.MenuID, query.PermissionID, query.PermissionID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (r * RequestPermissionRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from request_permission where request_permission_id in ?", ids).Error
}

func (r *RequestPermissionRepository) Save(requestPermission *model.RequestPermission) error {
	return r.BaseRepository.Create(requestPermission)
}

func (r *RequestPermissionRepository) Update(requestPermission *model.RequestPermission) error {
	return r.BaseRepository.Update(requestPermission)
}

