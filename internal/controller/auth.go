package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)


type IAuthController interface {
	GetToken() gin.HandlerFunc
	Register() gin.HandlerFunc
}

type AuthController struct {
	BaseController[model.TokenResponse]
	AuthService service.IAuthService
}

func NewAuthController() IAuthController {
	authService := service.NewAuthService()
	return &AuthController{AuthService: authService}
}

func (a AuthController) GetToken() gin.HandlerFunc {
	return a.ResponsePointer(a.AuthService.Authentication)
}

func (a AuthController) Register() gin.HandlerFunc {
	base := BaseController[dtos.UserResponse]{}
	return base.ResponsePointer(a.AuthService.Register)
}
