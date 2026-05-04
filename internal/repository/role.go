package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	IBaseRepository[model.Role, int]
	GetByRolename(roleName string) (*model.Role, error)
	GetAll() ([]model.Role, error)
	GetAllByCondition(query dtos.RoleFilter) ([]model.Role, int, error)
	Delete(ids []int) error
	GetRoleById(roleId int) (model.Role, error)
	Save(role *model.Role) error
	Update(role *model.Role) error
}

type RoleRepository struct {
	BaseRepository[model.Role, int]
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
		`select r.* 
	 from role r 
	 where (? is null or r.role_name like ?)
	 and (? is null OR r.created_at >= ?)
	 and (? is null OR r.created_at < ?)`,
		query.Page,
		query.Size,
		query.RoleName,
		query.RoleName,
		query.GetDateTo(),
		query.GetDateTo(),
		query.GetDateFrom(),
		query.GetDateFrom(),
	)
}
func (r *RoleRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from role where role_id in ?", ids).Error
}
func (r *RoleRepository) GetRoleById(roleId int) (model.Role, error) {
    var role model.Role
    err := r.DB.Where("role_id = ?", roleId).First(&role).Error
    if err != nil {
        return model.Role{}, err
    }
    return role, nil
}

func (r *RoleRepository) Save(role *model.Role) error {
	return r.BaseRepository.Create(role)
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.BaseRepository.Update(role)
}
