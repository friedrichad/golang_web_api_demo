package main

import (
	"github.com/friedrichad/golang_web_api_demo/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Server API
// @version 1.0
// @description JWT Auth Server
// @host localhost:8080
// @BasePath /

func main() {

	r := router.SetupRouter()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
