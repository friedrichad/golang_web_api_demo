package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
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
	BaseController[dtos.UserResponse]
	userService service.IUserService
}

func NewUserController() IUserController {
	userService := service.NewUserService()
	return &UserController{userService: userService}
}

func (controller *UserController) GetAllUsers() gin.HandlerFunc {
	return controller.ResponsePage(controller.userService.GetAllUsers)
}

func (controller *UserController) GetUserById() gin.HandlerFunc {
	return controller.ResponsePointer(controller.userService.GetUserByUuid)
}

func (controller *UserController) CreateUser() gin.HandlerFunc {
	return controller.ResponsePointer(controller.userService.CreateUser)
}

func (controller *UserController) UpdateUser() gin.HandlerFunc {
	return controller.ResponseSuccessOnly(controller.userService.UpdateUser)
}

func (controller *UserController) DeleteUser() gin.HandlerFunc {
	return controller.ResponseSuccessOnly(controller.userService.DeleteUser)
}

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
