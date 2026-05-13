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
	GetRoleById(roleId int) (*model.Role, error)
	Save(role *model.Role) error
	Update(role *model.Role) error

	// role_menu
	GetRoleMenus(roleId int) ([]model.RoleMenu, error)
	DeleteRoleMenus(roleId int) error
	CreateRoleMenus(roleMenus []model.RoleMenu) error

	// permissions
	GetAllPermissions() ([]model.Permission, error)
	GetPermissionById(id int) (*model.Permission, error)
	CreatePermission(permission *model.Permission) error
	UpdatePermission(permission *model.Permission) error
	DeletePermissions(ids []int) error

	// transaction methods
	CreateRoleTx(role *model.Role) (*model.Role, error)
	UpdateRoleTx(role *model.Role) error
	DeleteRoleTx(ids []int) error
	AssignRoleMenusTx(roleId int, menus []model.RoleMenu) error
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
func (r *RoleRepository) GetRoleById(roleId int) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("role_id = ?", roleId).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Save(role *model.Role) error {
	return r.BaseRepository.Create(role)
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.BaseRepository.Update(role)
}

func (r *RoleRepository) WithTx(tx *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: BaseRepository[model.Role, int]{Instance: tx},
		DB:             tx,
	}
}

// RoleMenu implementations
func (r *RoleRepository) GetRoleMenus(roleId int) ([]model.RoleMenu, error) {
	var roleMenus []model.RoleMenu
	err := r.DB.Where("role_id = ?", roleId).Find(&roleMenus).Error
	return roleMenus, err
}

func (r *RoleRepository) DeleteRoleMenus(roleId int) error {
	return r.DB.Where("role_id = ?", roleId).Delete(&model.RoleMenu{}).Error
}

func (r *RoleRepository) CreateRoleMenus(roleMenus []model.RoleMenu) error {
	if len(roleMenus) == 0 {
		return nil
	}
	return r.DB.Create(&roleMenus).Error
}

// Permission implementations
func (r *RoleRepository) GetAllPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.Find(&permissions).Error
	return permissions, err
}

func (r *RoleRepository) GetPermissionById(id int) (*model.Permission, error) {
	var permission model.Permission
	err := r.DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *RoleRepository) CreatePermission(permission *model.Permission) error {
	return r.DB.Create(permission).Error
}

func (r *RoleRepository) UpdatePermission(permission *model.Permission) error {
	return r.DB.Model(permission).Updates(permission).Error
}

func (r *RoleRepository) DeletePermissions(ids []int) error {
	return r.DB.Where("permission_id IN ?", ids).Delete(&model.Permission{}).Error
}

// CreateRoleTx handles transaction for role creation
func (r *RoleRepository) CreateRoleTx(role *model.Role) (*model.Role, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := r.WithTx(tx)
	err := roleRepoTx.Save(role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return role, nil
}

// UpdateRoleTx handles transaction for role update
func (r *RoleRepository) UpdateRoleTx(role *model.Role) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := r.WithTx(tx)
	err := roleRepoTx.Update(role)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteRoleTx handles transaction for role deletion
func (r *RoleRepository) DeleteRoleTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := r.WithTx(tx)
	err := roleRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// AssignRoleMenusTx handles transaction for assigning menus to role
func (r *RoleRepository) AssignRoleMenusTx(roleId int, menus []model.RoleMenu) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := r.WithTx(tx)

	// Delete existing role menus
	if err := roleRepoTx.DeleteRoleMenus(roleId); err != nil {
		tx.Rollback()
		return err
	}

	// Create new role menus
	if err := roleRepoTx.CreateRoleMenus(menus); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
