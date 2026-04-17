package v1handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	dtos "github.com/friedrichad/golang_web_api_demo/dtos"
	"github.com/friedrichad/golang_web_api_demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)+$`)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	userService := &service.UserService{DB: db}
	return &UserHandler{
		userService: userService,
	}
}
func (h *UserHandler) GetUser(ctx *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch users: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
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

	user, err := h.userService.GetUserByID(int32(userInt))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
func (h *UserHandler) PostUser(ctx *gin.Context) {
	var user dtos.UserResp

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	createdUser, err := h.userService.PostUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": createdUser,
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
