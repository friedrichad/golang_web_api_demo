package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IBinController interface {
	GetAllBins() gin.HandlerFunc
	GetBinById() gin.HandlerFunc
	CreateBin() gin.HandlerFunc
	UpdateBin() gin.HandlerFunc
	DeleteBin() gin.HandlerFunc
}

type BinController struct {
	BaseController[dtos.BinResponse]
	binService service.IBinService
}

func NewBinController() IBinController {
	binService := service.NewBinService()
	return &BinController{binService: binService}
}

func (c *BinController) GetAllBins() gin.HandlerFunc {
	return c.ResponsePage(c.binService.GetAllBins)
}

func (c *BinController) GetBinById() gin.HandlerFunc {
	return c.ResponsePointer(c.binService.GetBinByBinId)
}

func (c *BinController) CreateBin() gin.HandlerFunc {
	return c.ResponsePointer(c.binService.CreateBin)
}

func (c *BinController) UpdateBin() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.binService.UpdateBin)
}

func (c *BinController) DeleteBin() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.binService.DeleteBin)
}