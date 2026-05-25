package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IRequestController interface {
	GetAllRequests() gin.HandlerFunc
	GetRequestById() gin.HandlerFunc
	CreateRequest() gin.HandlerFunc
	UpdateRequest() gin.HandlerFunc
	DeleteRequest() gin.HandlerFunc
	ApprovalRequest() gin.HandlerFunc
	ConfirmRequest() gin.HandlerFunc
}

type RequestController struct {
	BaseController[model.RequestResponse]
	requestService service.IRequestService
}

func NewRequestController() IRequestController {
	return &RequestController{
		requestService: service.NewRequestService(),
	}
}

// GetAllRequests godoc
// @Summary Get all requests
// @Description Get a list of all requests with optional filtering and pagination
// @Tags requests
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.RequestResponse
// @Failure 400 {object} common.Error
// @Router /requests [get]
func (c *RequestController) GetAllRequests() gin.HandlerFunc {
	return c.ResponsePage(c.requestService.GetAllRequests)
}

// GetRequestById godoc
// @Summary Get a request by ID
// @Description Retrieve details of a specific request by its ID
// @Tags requests
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.RequestResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /requests/{id} [get]
func (c *RequestController) GetRequestById() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.GetRequestById)
}

// CreateRequest godoc
// @Summary Create a new request
// @Description Add a new request to the system
// @Tags requests
// @Accept json
// @Produce json
// @Param request body model.RequestCreate true "Request Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.RequestResponse
// @Failure 400 {object} common.Error
// @Router /requests [post]
func (c *RequestController) CreateRequest() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.CreateRequest)
}

// UpdateRequest godoc
// @Summary Update an existing request
// @Description Modify details of an existing request
// @Tags requests
// @Accept json
// @Produce json
// @Param request body model.RequestUpdate true "Request Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /requests [put]
func (c *RequestController) UpdateRequest() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.UpdateRequest)
}

// DeleteRequest godoc
// @Summary Delete requests
// @Description Delete one or multiple requests by their IDs
// @Tags requests
// @Accept json
// @Produce json
// @Param request body []int true "Array of Request IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /requests [delete]
func (c *RequestController) DeleteRequest() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.DeleteRequest)
}

// ApprovalRequest godoc
// @Summary Approve a request
// @Description Approve a specific request
// @Tags requests
// @Accept json
// @Produce json
// @Param request body model.ApprovalRequest true "Approval Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /requests/approval [post]
func (c *RequestController) ApprovalRequest() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.ApprovalRequest)
}

// ConfirmRequest godoc
// @Summary Confirm a request
// @Description Confirm a specific request
// @Tags requests
// @Accept json
// @Produce json
// @Param request body model.ConfirmRequest true "Confirmation Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /requests/confirm [post]
func (c *RequestController) ConfirmRequest() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.ConfirmRequest)
}
