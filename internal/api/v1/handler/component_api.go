package v1handler

import (
	"fmt"
	"net/http"
	"strconv"

	dtos "github.com/friedrichad/golang_web_api_demo/dtos"
	repository "github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComponentHandler struct {
	componentService repository.IComponentRepository
}

func NewComponentHandler(db *gorm.DB) *ComponentHandler {
	componentService := &repository.ComponentRepository{DB: db}
	return &ComponentHandler{
		componentService: componentService,
	}
}
func (h *ComponentHandler) GetComponent(ctx *gin.Context) {
	componentResps, err := h.componentService.GetComponents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch components: %v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, componentResps)
}
func (h *ComponentHandler) GetComponentByID(ctx *gin.Context) {
	componentID, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid component ID: %v", componentID),
		})
		return
	}
	componentResp, err := h.componentService.GetComponentByID(int32(componentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch component: %v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, componentResp)
}
func (h *ComponentHandler) CreateComponent(ctx *gin.Context) {
	var compoentReq dtos.ComponentRequest
	if err := ctx.BindJSON(&compoentReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}
	componentResp, err := h.componentService.CreateComponent(&compoentReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create component: %v", err),
		})
		return
	}
	ctx.JSON(http.StatusCreated, componentResp)
}
