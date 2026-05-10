package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	GetAllRoles(c *gin.Context) ([]dtos.RoleResponse, int, *common.Error)
	GetRoleById(c *gin.Context) (*dtos.RoleResponse, *common.Error)
	CreateRole(c *gin.Context) (*dtos.RoleResponse, *common.Error)
	UpdateRole(c *gin.Context) *common.Error
	DeleteRole(c *gin.Context) *common.Error

	// role_menu
	AssignRoleMenus(c *gin.Context) *common.Error
	GetRoleMenus(c *gin.Context) ([]int, *common.Error)

	// permission
	GetAllPermissions(c *gin.Context) ([]dtos.PermissionDTO, *common.Error)
	GetPermissionById(c *gin.Context) (*dtos.PermissionDTO, *common.Error)
	CreatePermission(c *gin.Context) (*dtos.PermissionDTO, *common.Error)
	UpdatePermission(c *gin.Context) *common.Error
	DeletePermissions(c *gin.Context) *common.Error
}

type RoleService struct {
	roleRepo repository.IRoleRepository
}

var roleService IRoleService

func NewRoleService() IRoleService {
	if roleService == nil {
		roleService = &RoleService{
			roleRepo: repository.NewRoleRepository(),
		}
	}
	return roleService
}

func (s *RoleService) GetAllRoles(c *gin.Context) ([]dtos.RoleResponse, int, *common.Error) {
	var query dtos.RoleFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	roles, total, err := s.roleRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	roleResponses := make([]dtos.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = modelToRoleResponse(&role)
	}

	return roleResponses, total, nil
}

func (s *RoleService) GetRoleById(c *gin.Context) (*dtos.RoleResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	roleId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	role, err := s.roleRepo.GetRoleById(int(roleId))
	if err != nil {
		return nil, common.NotFound
	}

	if role == nil {
		return nil, &common.Error{Code: "404", Message: "Quyền không tồn tại"}
	}

	roleResponse := modelToRoleResponse(role)
	return &roleResponse, nil
}

func (s *RoleService) CreateRole(c *gin.Context) (*dtos.RoleResponse, *common.Error) {
	var req dtos.RoleCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return nil, &common.Error{Code: "400", Message: err.Error()}
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := s.roleRepo.(*repository.RoleRepository).WithTx(tx)

	role := &model.Role{
		RoleName:    req.RoleName,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	err := roleRepoTx.Save(role)
	if err != nil {
		tx.Rollback()
		return nil, common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, common.SystemError
	}

	roleResponse := modelToRoleResponse(role)
	return &roleResponse, nil
}

func (s *RoleService) UpdateRole(c *gin.Context) *common.Error {
	var req dtos.RoleUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	if err := req.Verify(); err != nil {
		return &common.Error{Code: "400", Message: err.Error()}
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := s.roleRepo.(*repository.RoleRepository).WithTx(tx)

	role, err := roleRepoTx.GetRoleById(int(req.RoleID))
	if err != nil {
		tx.Rollback()
		return common.NotFound
	}
	if role == nil {
		tx.Rollback()
		return &common.Error{Code: "404", Message: "Quyền không tồn tại"}
	}

	if req.RoleName != "" {
		role.RoleName = req.RoleName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	role.UpdatedBy = int(req.UpdatedBy)
	role.UpdatedAt = time.Now()

	err = roleRepoTx.Update(role)
	if err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func (s *RoleService) DeleteRole(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		roleId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(roleId)
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := s.roleRepo.(*repository.RoleRepository).WithTx(tx)

	err := roleRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func modelToRoleResponse(role *model.Role) dtos.RoleResponse {
	return dtos.RoleResponse{
		RoleID:      int(role.RoleID),
		RoleName:    role.RoleName,
		Description: role.Description,
		CreatedBy:   role.CreatedBy,
		CreatedAt:   role.CreatedAt,
		UpdatedBy:   role.UpdatedBy,
		UpdatedAt:   role.UpdatedAt,
	}
}

// role_menu
func (s *RoleService) AssignRoleMenus(c *gin.Context) *common.Error {
	var req dtos.RoleMenuAssign
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	roleRepoTx := s.roleRepo.(*repository.RoleRepository).WithTx(tx)

	if err := roleRepoTx.DeleteRoleMenus(req.RoleID); err != nil {
		tx.Rollback()
		return common.SystemError
	}

	roleMenus := make([]model.RoleMenu, len(req.MenuIDs))
	for i, menuId := range req.MenuIDs {
		roleMenus[i] = model.RoleMenu{
			RoleID: req.RoleID,
			MenuID: menuId,
		}
	}

	if err := roleRepoTx.CreateRoleMenus(roleMenus); err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func (s *RoleService) GetRoleMenus(c *gin.Context) ([]int, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	roleId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	roleMenus, err := s.roleRepo.GetRoleMenus(int(roleId))
	if err != nil {
		return nil, common.SystemError
	}

	menuIds := make([]int, len(roleMenus))
	for i, rm := range roleMenus {
		menuIds[i] = rm.MenuID
	}
	return menuIds, nil
}

// permission
func (s *RoleService) GetAllPermissions(c *gin.Context) ([]dtos.PermissionDTO, *common.Error) {
	permissions, err := s.roleRepo.GetAllPermissions()
	if err != nil {
		return nil, common.SystemError
	}

	permissionDTOs := make([]dtos.PermissionDTO, len(permissions))
	for i, p := range permissions {
		permissionDTOs[i] = dtos.PermissionDTO{
			PermissionID:   p.PermissionID,
			PermissionName: p.PermissionName,
		}
	}
	return permissionDTOs, nil
}

func (s *RoleService) GetPermissionById(c *gin.Context) (*dtos.PermissionDTO, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	permId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	permission, err := s.roleRepo.GetPermissionById(int(permId))
	if err != nil {
		return nil, common.NotFound
	}

	return &dtos.PermissionDTO{
		PermissionID:   permission.PermissionID,
		PermissionName: permission.PermissionName,
	}, nil
}

func (s *RoleService) CreatePermission(c *gin.Context) (*dtos.PermissionDTO, *common.Error) {
	var req dtos.PermissionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	permission := &model.Permission{
		PermissionName: req.PermissionName,
	}

	if err := s.roleRepo.CreatePermission(permission); err != nil {
		return nil, common.SystemError
	}

	return &dtos.PermissionDTO{
		PermissionID:   permission.PermissionID,
		PermissionName: permission.PermissionName,
	}, nil
}

func (s *RoleService) UpdatePermission(c *gin.Context) *common.Error {
	var req dtos.PermissionUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	permission, err := s.roleRepo.GetPermissionById(req.PermissionID)
	if err != nil || permission == nil {
		return common.NotFound
	}

	permission.PermissionName = req.PermissionName

	if err := s.roleRepo.UpdatePermission(permission); err != nil {
		return common.SystemError
	}

	return nil
}

func (s *RoleService) DeletePermissions(c *gin.Context) *common.Error {
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		permId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(permId)
	}

	if err := s.roleRepo.DeletePermissions(ids); err != nil {
		return common.SystemError
	}

	return nil
}
