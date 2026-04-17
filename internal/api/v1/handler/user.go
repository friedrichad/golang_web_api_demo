package v1handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)+$`)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}
func (h *UserHandler) GetUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "GetUsers method",
	})
}
func (h *UserHandler) GetUserSlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if !slugRegex.MatchString(slug) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user slug",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "GetUsers method",
	})
}
func (h *UserHandler) GetUserById(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	userInt, err := strconv.Atoi(user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	if userInt < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID must be a positive integer",
		})
		return
	}
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
