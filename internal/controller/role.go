package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IRoleController interface {
	GetAllRoles() gin.HandlerFunc
	GetRoleById() gin.HandlerFunc
	CreateRole() gin.HandlerFunc
	UpdateRole() gin.HandlerFunc
	DeleteRole() gin.HandlerFunc
}

type RoleController struct {
	BaseController[model.RoleResponse]
	roleService service.IRoleService
}

func NewRoleController() IRoleController {
	return &RoleController{
		roleService: service.NewRoleService(),
	}
}

// GetAllRoles godoc
// @Summary Get all roles
// @Description Get a list of all roles with optional filtering and pagination
// @Tags roles
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.RoleResponse
// @Failure 400 {object} common.Error
// @Router /roles [get]
func (c *RoleController) GetAllRoles() gin.HandlerFunc {
	return c.ResponsePage(c.roleService.GetAllRoles)
}

// GetRoleById godoc
// @Summary Get a role by ID
// @Description Retrieve details of a specific role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.RoleResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /roles/{id} [get]
func (c *RoleController) GetRoleById() gin.HandlerFunc {
	return c.ResponsePointer(c.roleService.GetRoleById)
}

// CreateRole godoc
// @Summary Create a new role
// @Description Add a new role to the system
// @Tags roles
// @Accept json
// @Produce json
// @Param request body model.RoleCreate true "Role Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.RoleResponse
// @Failure 400 {object} common.Error
// @Router /roles [post]
func (c *RoleController) CreateRole() gin.HandlerFunc {
	return c.ResponsePointer(c.roleService.CreateRole)
}

// UpdateRole godoc
// @Summary Update an existing role
// @Description Modify details of an existing role
// @Tags roles
// @Accept json
// @Produce json
// @Param request body model.RoleUpdate true "Role Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /roles [put]
func (c *RoleController) UpdateRole() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.roleService.UpdateRole)
}

// DeleteRole godoc
// @Summary Delete roles
// @Description Delete one or multiple roles by their IDs
// @Tags roles
// @Accept json
// @Produce json
// @Param request body []int true "Array of Role IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /roles [delete]
func (c *RoleController) DeleteRole() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.roleService.DeleteRole)
}
