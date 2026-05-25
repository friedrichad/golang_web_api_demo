package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
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
	BaseController[model.BinResponse]
	binService service.IBinService
}

func NewBinController() IBinController {
	binService := service.NewBinService()
	return &BinController{binService: binService}
}

// GetAllBins godoc
// @Summary Get all bins
// @Description Get a list of all bins with optional filtering and pagination
// @Tags bins
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.BinResponse
// @Failure 400 {object} common.Error
// @Router /bins [get]
func (c *BinController) GetAllBins() gin.HandlerFunc {
	return c.ResponsePage(c.binService.GetAllBins)
}

// GetBinById godoc
// @Summary Get a bin by ID
// @Description Retrieve details of a specific bin by its ID
// @Tags bins
// @Accept json
// @Produce json
// @Param id path int true "Bin ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.BinResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /bins/{id} [get]
func (c *BinController) GetBinById() gin.HandlerFunc {
	return c.ResponsePointer(c.binService.GetBinByBinId)
}

// CreateBin godoc
// @Summary Create a new bin
// @Description Add a new bin to the system
// @Tags bins
// @Accept json
// @Produce json
// @Param request body model.BinCreate true "Bin Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.BinResponse
// @Failure 400 {object} common.Error
// @Router /bins [post]
func (c *BinController) CreateBin() gin.HandlerFunc {
	return c.ResponsePointer(c.binService.CreateBin)
}

// UpdateBin godoc
// @Summary Update an existing bin
// @Description Modify details of an existing bin
// @Tags bins
// @Accept json
// @Produce json
// @Param request body model.BinUpdate true "Bin Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /bins [put]
func (c *BinController) UpdateBin() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.binService.UpdateBin)
}

// DeleteBin godoc
// @Summary Delete bins
// @Description Delete one or multiple bins by their IDs
// @Tags bins
// @Accept json
// @Produce json
// @Param request body []int true "Array of Bin IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /bins [delete]
func (c *BinController) DeleteBin() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.binService.DeleteBin)
}
