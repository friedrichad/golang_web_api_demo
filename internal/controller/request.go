package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IRequestController interface {
	GetAllRequests() gin.HandlerFunc
	GetRequestById() gin.HandlerFunc
	CreateRequest() gin.HandlerFunc
	UpdateRequest() gin.HandlerFunc
	DeleteRequest() gin.HandlerFunc
}

type RequestController struct {
	BaseController[dtos.RequestResponse]
	requestService service.IRequestService
}

func NewRequestController() IRequestController {
	return &RequestController{
		requestService: service.NewRequestService(),
	}
}

func (c *RequestController) GetAllRequests() gin.HandlerFunc {
	return c.ResponsePage(c.requestService.GetAllRequests)
}

func (c *RequestController) GetRequestById() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.GetRequestById)
}

func (c *RequestController) CreateRequest() gin.HandlerFunc {
	return c.ResponsePointer(c.requestService.CreateRequest)
}

func (c *RequestController) UpdateRequest() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.UpdateRequest)
}

func (c *RequestController) DeleteRequest() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.requestService.DeleteRequest)
}
