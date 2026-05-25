package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	GetAllUsers() gin.HandlerFunc
	GetUserById() gin.HandlerFunc
	CreateUser() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	GetUserAuthorities() gin.HandlerFunc
}

type UserController struct {
	BaseController[model.UserResponse]
	userService service.IUserService
}

func NewUserController() IUserController {
	userService := service.NewUserService()
	return &UserController{userService: userService}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users with optional filtering and pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.UserResponse
// @Failure 400 {object} common.Error
// @Router /users [get]
func (controller *UserController) GetAllUsers() gin.HandlerFunc {
	return controller.ResponsePage(controller.userService.GetAllUsers)
}

// GetUserById godoc
// @Summary Get a user by ID
// @Description Retrieve details of a specific user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Security ApiKeyAuth
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /users/{id} [get]
func (controller *UserController) GetUserById() gin.HandlerFunc {
	return controller.ResponsePointer(controller.userService.GetUserByUuid)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Add a new user to the system
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.UserCreate true "User Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} common.Error
// @Router /users [post]
func (controller *UserController) CreateUser() gin.HandlerFunc {
	return controller.ResponsePointer(controller.userService.CreateUser)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Modify details of an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.UserUpdate true "User Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /users [put]
func (controller *UserController) UpdateUser() gin.HandlerFunc {
	return controller.ResponseSuccessOnly(controller.userService.UpdateUser)
}

// DeleteUser godoc
// @Summary Delete users
// @Description Delete one or multiple users by their UUIDs
// @Tags users
// @Accept json
// @Produce json
// @Param request body []string true "Array of User UUIDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /users [delete]
func (controller *UserController) DeleteUser() gin.HandlerFunc {
	return controller.ResponseSuccessOnly(controller.userService.DeleteUser)
}

// GetUserAuthorities godoc
// @Summary Get user authorities
// @Description Retrieve all authorities assigned to a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Security ApiKeyAuth
// @Success 200 {object} []string
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /users/{id}/authorities [get]
func (controller *UserController) GetUserAuthorities() gin.HandlerFunc {
	return func(g *gin.Context) {
		authorities, err := controller.userService.GetUserAuthorities(g)
		if err != nil {
			controller.Error(g, err, nil)
			return
		}
		controller.Success(g, authorities)
	}
}
