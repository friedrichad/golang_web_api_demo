package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IMenuPermissionRepository interface {
	GetAllMenuPermissions(query *model.MenuPermissionFilter) ([]model.MenuPermission, int, error)
	GetAllMenuPermissionsByRestricted(isRestricted int) ([]string, error)
	GetMenuPermissionById(id int) (*model.MenuPermission, error)
	AddMenuPermission(menuPermission *model.MenuPermission) error
	UpdateMenuPermission(menuPermission *model.MenuPermission) error
	DeleteMenuPermissions(menuPermissionId int) error
}

type MenuPermissionRepository struct {
	BaseRepository[model.MenuPermission, int]
	DB *gorm.DB
}

var menuPermissionRepository *MenuPermissionRepository

func NewMenuPermissionRepository() IMenuPermissionRepository {
	if menuPermissionRepository == nil {
		menuPermissionRepository = &MenuPermissionRepository{DB: db.Instance}
		menuPermissionRepository.SetInstance(db.Instance)
	}
	return menuPermissionRepository
}

func (m *MenuPermissionRepository) GetAllMenuPermissions(query *model.MenuPermissionFilter) ([]model.MenuPermission, int, error) {
	return m.GetPage(
		"SELECT mp.menu_permission_id, "+
			"m.menu_name, "+
			"p.permission_name, "+
			"CONCAT(m.menu_name, ':', p.permission_name) AS scope, "+
			"mp.is_restricted "+
			"FROM menu_permission mp "+
			"JOIN menu m ON mp.menu_id = m.menu_id "+
			"JOIN permissions p ON mp.permission_id = p.permission_id "+
			"WHERE (? IS NULL OR m.menu_name LIKE CONCAT('%', ?, '%')) "+
			"AND (? IS NULL OR p.permission_name LIKE CONCAT('%', ?, '%')) "+
			"AND (? IS NULL OR mp.is_restricted = ?) "+
			"AND (? IS NULL OR mp.created_at >= ?) "+
			"AND (? IS NULL OR mp.created_at <= ?) ", query.Page, query.Size, query.MenuPermissionName, query.MenuPermissionName, query.PermissionName, query.PermissionName, query.IsRestricted, query.IsRestricted, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo(),
	)
}

func (m *MenuPermissionRepository) GetMenuPermissionById(id int) (*model.MenuPermission, error) {
	return m.GetById(id)
}

func (m *MenuPermissionRepository) AddMenuPermission(menuPermission *model.MenuPermission) error {
	return m.BaseRepository.Create(menuPermission)
}

func (m *MenuPermissionRepository) UpdateMenuPermission(menuPermission *model.MenuPermission) error {
	return m.BaseRepository.Update(menuPermission)
}

func (m *MenuPermissionRepository) DeleteMenuPermissions(menuPermissionId int) error {
	return m.DB.Exec("DELETE FROM menu_permission WHERE menu_permission_id = ?", menuPermissionId).Error
}

func (m *MenuPermissionRepository) GetAllMenuPermissionsByRestricted(isRestricted int) ([]string, error) {
	var scopes []string
	err := m.DB.Raw(
		"SELECT CONCAT(m.menu_name, ':', p.permission_name) AS scope "+
			"FROM menu_permission mp "+
			"JOIN menu m "+
			"ON mp.menu_id = m.menu_id "+
			"JOIN permissions p "+
			"ON mp.permission_id = p.permission_id "+
			"WHERE ? IS NULL OR mp.is_restricted = ? "+
			"ORDER BY mp.menu_permission_id DESC",
		isRestricted, isRestricted,
	).Scan(&scopes).Error
	return scopes, err
}
