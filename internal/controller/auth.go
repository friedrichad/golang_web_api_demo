package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	GetToken() gin.HandlerFunc
	Register() gin.HandlerFunc
	Logout() gin.HandlerFunc
}

type AuthController struct {
	BaseController[model.TokenResponse]
	AuthService service.IAuthService
}

func NewAuthController() IAuthController {
	authService := service.NewAuthService()
	return &AuthController{AuthService: authService}
}

// GetToken godoc
// @Summary Get authentication token
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant Type (password or refresh_token)"
// @Param username formData string false "Username"
// @Param password formData string false "Password"
// @Param refresh_token formData string false "Refresh Token"
// @Success 200 {object} model.TokenResponse
// @Router /auth/login [post]
func (a AuthController) GetToken() gin.HandlerFunc {
	return a.ResponsePointer(a.AuthService.Authentication)
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration Info"
// @Success 200 {object} model.UserResponse
// @Router /auth/register [post]
func (a AuthController) Register() gin.HandlerFunc {
	base := BaseController[model.UserResponse]{}
	return base.ResponsePointer(a.AuthService.Register)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate the current JWT token
// @Tags auth
// @Produce json
// @Success 200 "Success"
// @Security ApiKeyAuth
// @Router /auth/logout [post]
func (a AuthController) Logout() gin.HandlerFunc {
	return a.ResponseSuccessOnly(a.AuthService.Logout)
}
