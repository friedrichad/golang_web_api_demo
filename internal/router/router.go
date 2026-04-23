package router

import (
	// "github.com/friedrichad/golang_web_api_demo/internal/controller"
	// "github.com/friedrichad/golang_web_api_demo/internal/configs/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	configCors(router)
	return router
}

func configCors(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:  viper.GetStringSlice("security.cors"),
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))
}

func initOtherRouter(router *gin.Engine) {
	// router.GET("role-menu", middleware.BearerAuthenticator(), controller.NewRoleMenuController().GetByRole())
	// router.GET("user-role", middleware.BearerAuthenticator(), controller.NewUserRoleController().GetByUser())
	// router.POST("upload", middleware.BearerAuthenticator(), controller.NewUploadController().UploadFile())
}
