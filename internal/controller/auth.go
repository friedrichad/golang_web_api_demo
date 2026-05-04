package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)


type IStorageController interface {
	GetToken() gin.HandlerFunc
}

type AuthController struct {
	BaseController[model.TokenResponse]
	AuthService service.IAuthService
}

func NewAuthController() IStorageController {
	authService := service.NewAuthService()
	return &AuthController{AuthService: authService}
}

func (a AuthController) GetToken() gin.HandlerFunc {
	return a.ResponsePointer(a.AuthService.Authentication)
}
