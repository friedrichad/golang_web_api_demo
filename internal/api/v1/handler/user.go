package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}
func (h *UserHandler) GetUserById(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	ctx.JSON(http.StatusOK, gin.H{
		"data":    "GetUserById method",
		"user_id": user_id,
	})
}
func (h *UserHandler) PostUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "PostUser method",
	})
}
func (h *UserHandler) PutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "PutUser method",
	})
}
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "DeleteUser method",
	})
}
