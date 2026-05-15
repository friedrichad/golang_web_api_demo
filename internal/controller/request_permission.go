package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IRequestPermissionController interface {
	GetAllPermissions() gin.HandlerFunc
	CreatePermission() gin.HandlerFunc
	UpdatePermission() gin.HandlerFunc
	DeletePermission() gin.HandlerFunc
	ApprovalPermission() gin.HandlerFunc
}

type RequestPermissionController struct {
	BaseController[model.RequestPermissionResponse]
	requestPermissionService service.IRequestPermissionService
}

func NewRequestPermissionController() IRequestPermissionController {
	return &RequestPermissionController{
		requestPermissionService: service.NewRequestPermissionService(),
	}
}

func (c *RequestPermissionController) GetAllPermissions() gin.HandlerFunc {
	return c.ResponsePage(c.requestPermissionService.GetAllPermissionByCondition)
}

func (c *RequestPermissionController) CreatePermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Save)
}

func (c *RequestPermissionController) UpdatePermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Update)
}

func (c *RequestPermissionController) DeletePermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Delete)
}

func (c *RequestPermissionController) ApprovalPermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Approval)
}
