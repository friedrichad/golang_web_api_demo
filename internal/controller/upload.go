package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IUploadController interface {
	UploadBase64() gin.HandlerFunc
	UploadMultipart() gin.HandlerFunc
	UploadMultiple() gin.HandlerFunc
}

type UploadController struct {
	BaseController[dtos.UploadResponse]
	uploadService service.IUploadService
}

func NewUploadController() IUploadController {
	uploadService := service.NewUploadService()
	return &UploadController{uploadService: uploadService}
}

func (c *UploadController) UploadBase64() gin.HandlerFunc {
	return c.ResponsePointer(c.uploadService.UploadBase64)
}

func (c *UploadController) UploadMultipart() gin.HandlerFunc {
	return c.ResponsePointer(c.uploadService.UploadMultipart)
}
func (c *UploadController) UploadMultiple() gin.HandlerFunc {
	return c.ResponsePage(c.uploadService.UploadMultiple)
}