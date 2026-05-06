package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
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
	BaseController[dtos.RequestDetailResponse]
	requestService service.IRequestService
}

func NewRequestDetailController() IRequestDetailController {
	return &RequestDetailController{
		requestService: service.NewRequestService(),
	}
}

func (c *RequestDetailController) GetAllRequestDetails() gin.HandlerFunc {
	return c.ResponsePage(c.requestService.GetAllRequestDetails)
}

func (c *RequestDetailController) GetRequestDetailById() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.GetRequestDetailById)
}

func (c *RequestDetailController) CreateRequestDetail() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.CreateRequestDetail)
}

func (c *RequestDetailController) UpdateRequestDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.UpdateRequestDetail)
}

func (c *RequestDetailController) DeleteRequestDetail() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.DeleteRequestDetail)
}
