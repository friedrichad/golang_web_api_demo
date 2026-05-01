package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
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
	BaseController[dtos.RoleResponse]
	roleService service.IRoleService
}

func NewRoleController() IRoleController {
	return &RoleController{
		roleService: service.NewRoleService(),
	}
}

func (c *RoleController) GetAllRoles() gin.HandlerFunc {
	return c.ResponsePage(c.roleService.GetAllRoles)
}

func (c *RoleController) GetRoleById() gin.HandlerFunc {
	return c.ResponsePointer(c.roleService.GetRoleById)
}

func (c *RoleController) CreateRole() gin.HandlerFunc {
	return c.ResponsePointer(c.roleService.CreateRole)
}

func (c *RoleController) UpdateRole() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.roleService.UpdateRole)
}

func (c *RoleController) DeleteRole() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.roleService.DeleteRole)
}
