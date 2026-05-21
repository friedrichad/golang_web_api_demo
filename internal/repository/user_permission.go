package repository

import (
	"log"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IUserPermissionRepository interface {
	Save(userPermission *model.UserPermission) error
	SaveBatch(userPermissions []model.UserPermission) error
	Delete(userPermissionIds []int) error
}

type UserPermissionRepository struct {
	BaseRepository[model.UserPermission, int]
	DB *gorm.DB
}

var userPermissionRepository IUserPermissionRepository

func NewUserPermissionRepository() IUserPermissionRepository {
	if userPermissionRepository == nil {
		userPermissionRepository = &UserPermissionRepository{DB: db.Instance}
	}
	return userPermissionRepository
}
func (r *UserPermissionRepository) Save(userPermission *model.UserPermission) error {
	err := r.DB.Create(userPermission).Error
	if err != nil {
		log.Print("Lỗi khi lưu user permission: ", err)
		return err
	}
	return nil
}

// SaveBatch creates multiple user permissions in a single batch query
func (r *UserPermissionRepository) SaveBatch(userPermissions []model.UserPermission) error {
	if len(userPermissions) == 0 {
		return nil
	}
	err := r.DB.CreateInBatches(userPermissions, 100).Error
	if err != nil {
		log.Print("Lỗi khi lưu batch user permissions: ", err)
		return err
	}
	return nil
}

func (r *UserPermissionRepository) Delete(userPermissionIds []int) error {
	err := r.DB.Where("menu_permission_id IN ?", userPermissionIds).Delete(&model.UserPermission{}).Error
	if err != nil {
		log.Print("Lỗi khi xóa user permission: ", err)
		return err
	}
	return nil
}