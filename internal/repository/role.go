package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	IBaseRepository[model.Role, string]
	GetByRolename(roleName string) (*model.Role, error)
	GetAll() ([]model.Role, error)
	GetAllByCondition(query dtos.RoleFilter) ([]model.Role, int, error)
	Delete(ids []string) error
	GetRoleByUuid(uuid string) ([]model.Role, error)
	Save(role *model.Role) error
	Update(role *model.Role) error
}

type RoleRepository struct {
	BaseRepository[model.Role, string]
	DB *gorm.DB
}

var roleRepository IRoleRepository

func NewRoleRepository() IRoleRepository {
	if roleRepository == nil {
		roleRepository = &RoleRepository{DB: db.Instance}
		roleRepository.SetInstance(db.Instance)
	}
	return roleRepository
}
func (r *RoleRepository) GetByRolename(roleName string) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("role_name = ?", roleName).First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &role, err
}

func (r *RoleRepository) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.Order("created_at").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetAllByCondition(query dtos.RoleFilter) ([]model.Role, int, error) {
	return r.GetPage(
		`SELECT r.* 
	 FROM role r 
	 WHERE r.role_name LIKE ?
	 AND (? IS NULL OR r.create_at >= ?)
	 AND (? IS NULL OR r.create_at < ?)`,
		query.Page,
		query.Size,
		"%"+query.RoleName+"%",
		query.GetDateTo(),
		query.GetDateTo(),
		query.GetDateFrom(),
		query.GetDateFrom(),
	)
}
func (r *RoleRepository) Delete(ids []string) error {
	return r.DB.Exec("delete r from r where r.role_id in ?", ids).Error
}
func (r *RoleRepository) GetRoleByUuid(uuid string) ([]model.Role, error) {
	var roles []model.Role

	err := r.DB.Raw(`
		SELECT r.*
		FROM role AS r
		JOIN user_role ur ON ur.role_id = r.role_id
		JOIN user u ON ur.user_uuid = u.user_uuid
		WHERE u.user_uuid = ?
	`, uuid).Scan(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}
func (r *RoleRepository) Save(role *model.Role) error {
	return r.BaseRepository.Create(role)
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.BaseRepository.Update(role)
}
