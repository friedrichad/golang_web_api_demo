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

// GetAllPermissions godoc
// @Summary Get all request permissions
// @Description Get a list of all request permissions with optional filtering and pagination
// @Tags request-permissions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.RequestPermissionResponse
// @Failure 400 {object} common.Error
// @Router /request-permissions [get]
func (c *RequestPermissionController) GetAllPermissions() gin.HandlerFunc {
	return c.ResponsePage(c.requestPermissionService.GetAllPermissionByCondition)
}

// CreatePermission godoc
// @Summary Create a new request permission
// @Description Add a new request permission to the system
// @Tags request-permissions
// @Accept json
// @Produce json
// @Param request body model.RequestPermissionCreate true "Permission Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /request-permissions [post]
func (c *RequestPermissionController) CreatePermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Save)
}

// UpdatePermission godoc
// @Summary Update an existing request permission
// @Description Modify details of an existing request permission
// @Tags request-permissions
// @Accept json
// @Produce json
// @Param request body model.RequestPermissionUpdate true "Permission Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /request-permissions [put]
func (c *RequestPermissionController) UpdatePermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Update)
}

// DeletePermission godoc
// @Summary Delete request permissions
// @Description Delete one or multiple request permissions by their IDs
// @Tags request-permissions
// @Accept json
// @Produce json
// @Param request body []int true "Array of Permission IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /request-permissions [delete]
func (c *RequestPermissionController) DeletePermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Delete)
}

// ApprovalPermission godoc
// @Summary Approve a request permission
// @Description Approve a specific request permission
// @Tags request-permissions
// @Accept json
// @Produce json
// @Param request body model.RequestPermissionUpdate true "Approval Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /request-permissions/approval [post]
func (c *RequestPermissionController) ApprovalPermission() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestPermissionService.Approval)
}
