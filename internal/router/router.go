package router

import (
	"github.com/friedrichad/golang_web_api_demo/internal/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

)

func InitRouter() *gin.Engine {
	router := gin.Default()
	configCors(router)
	// initOtherRouter(router)
	initUserRouter(router)
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



func initUserRouter(router *gin.Engine) {
	userController := controller.NewUserController()
	userGroup := router.Group("/users")
	{
		userGroup.GET("", userController.GetAllUsers)
		userGroup.GET("/:id", userController.GetUserById)
		userGroup.POST("", userController.CreateUser)
		userGroup.PUT("", userController.UpdateUser)
		userGroup.DELETE("", userController.DeleteUser)
		userGroup.GET("/:id/authorities", userController.GetUserAuthorities)
	}
}
