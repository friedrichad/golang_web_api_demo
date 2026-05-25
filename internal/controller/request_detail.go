package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IRequestDetailController interface {
	GetAllRequestDetails() gin.HandlerFunc
	GetRequestDetailById() gin.HandlerFunc
	CreateRequestDetail() gin.HandlerFunc
	UpdateRequestDetail() gin.HandlerFunc
	DeleteRequestDetail() gin.HandlerFunc
}

type RequestDetailController struct {
	BaseController[model.RequestDetailResponse]
	requestService service.IRequestService
}

func NewRequestDetailController() IRequestDetailController {
	return &RequestDetailController{
		requestService: service.NewRequestService(),
	}
}

// GetAllRequestDetails godoc
// @Summary Get all request details
// @Description Get a list of all request details with optional filtering and pagination
// @Tags request-details
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.RequestDetailResponse
// @Failure 400 {object} common.Error
// @Router /request-details [get]
func (c *RequestDetailController) GetAllRequestDetails() gin.HandlerFunc {
	return c.ResponsePage(c.requestService.GetAllRequestDetails)
}

// GetRequestDetailById godoc
// @Summary Get a request detail by ID
// @Description Retrieve details of a specific request detail by its ID
// @Tags request-details
// @Accept json
// @Produce json
// @Param id path int true "Request Detail ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.RequestDetailResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /request-details/{id} [get]
func (c *RequestDetailController) GetRequestDetailById() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.GetRequestDetailById)
}

// CreateRequestDetail godoc
// @Summary Create a new request detail
// @Description Add a new request detail to a request
// @Tags request-details
// @Accept json
// @Produce json
// @Param request body model.RequestDetailCreate true "Request Detail Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.RequestDetailResponse
// @Failure 400 {object} common.Error
// @Router /request-details [post]
func (c *RequestDetailController) CreateRequestDetail() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.CreateRequestDetail)
}

// UpdateRequestDetail godoc
// @Summary Update an existing request detail
// @Description Modify details of an existing request detail
// @Tags request-details
// @Accept json
// @Produce json
// @Param request body model.RequestDetailUpdate true "Request Detail Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /request-details [put]
func (c *RequestDetailController) UpdateRequestDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.UpdateRequestDetail)
}

// DeleteRequestDetail godoc
// @Summary Delete request details
// @Description Delete one or multiple request details by their IDs
// @Tags request-details
// @Accept json
// @Produce json
// @Param request body []int true "Array of Request Detail IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /request-details [delete]
func (c *RequestDetailController) DeleteRequestDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.DeleteRequestDetail)
}
