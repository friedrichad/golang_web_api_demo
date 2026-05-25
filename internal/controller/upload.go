package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IUploadController interface {
	UploadBase64() gin.HandlerFunc
	UploadMultipart() gin.HandlerFunc
	UploadMultiple() gin.HandlerFunc
}

type UploadController struct {
	BaseController[model.UploadResponse]
	uploadService service.IUploadService
}

func NewUploadController() IUploadController {
	uploadService := service.NewUploadService()
	return &UploadController{uploadService: uploadService}
}

// UploadBase64 godoc
// @Summary Upload file as Base64
// @Description Upload a file encoded in Base64 format
// @Tags uploads
// @Accept json
// @Produce json
// @Param request body model.UploadBase64Request true "Base64 file data"
// @Security ApiKeyAuth
// @Success 200 {object} model.UploadResponse
// @Failure 400 {object} common.Error
// @Router /uploads/base64 [post]
func (c *UploadController) UploadBase64() gin.HandlerFunc {
	return c.ResponsePointer(c.uploadService.UploadBase64)
}

// UploadMultipart godoc
// @Summary Upload file as multipart form data
// @Description Upload a single file as multipart form data
// @Tags uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Security ApiKeyAuth
// @Success 200 {object} model.UploadResponse
// @Failure 400 {object} common.Error
// @Router /uploads/multipart [post]
func (c *UploadController) UploadMultipart() gin.HandlerFunc {
	return c.ResponsePointer(c.uploadService.UploadMultipart)
}

// UploadMultiple godoc
// @Summary Upload multiple files
// @Description Upload multiple files at once
// @Tags uploads
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload"
// @Security ApiKeyAuth
// @Success 200 {object} []model.UploadResponse
// @Failure 400 {object} common.Error
// @Router /uploads/multiple [post]
func (c *UploadController) UploadMultiple() gin.HandlerFunc {
	return c.ResponsePage(c.uploadService.UploadMultiple)
}
