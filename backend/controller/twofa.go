package controller

import (
	"github.com/friedrichad/golang_web_api_demo/backend/model"
	"github.com/friedrichad/golang_web_api_demo/backend/service"
	"github.com/gin-gonic/gin"
)

type TwoFAController struct {
	service service.IAuthService
}

func NewTwoFAController() *TwoFAController {
	return &TwoFAController{
		service: service.NewAuthService(),
	}
}

func (t *TwoFAController) Generate2FA() gin.HandlerFunc {
	base := BaseController[model.TwoFASetupResponse]{}
	return base.ResponsePointer(t.service.Setup2FA)
}

func (t *TwoFAController) VerifySetup2FA() gin.HandlerFunc {
	base := BaseController[model.TwoFAVerifyResponse]{}
	return base.ResponsePointer(t.service.VerifySetup2FA)
}

func (t *TwoFAController) Verify2FA() gin.HandlerFunc {
	base := BaseController[model.TokenResponse]{}
	return base.ResponsePointer(t.service.Verify2FA)
}
