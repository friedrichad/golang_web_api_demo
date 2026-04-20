package main

import (
	db "github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	v1handler "github.com/friedrichad/golang_web_api_demo/internal/api/v1/handler"
	config "github.com/friedrichad/golang_web_api_demo/internal/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("internal/configs/config.json")
	if err != nil {
		panic(err)
	}

	db.ConnectDB(cfg)

	r := gin.Default()
	v1 := r.Group("/api/v1")

	// Khởi tạo UserHandler với DB đã connect
	userHandler := v1handler.NewUserHandler(db.DB)
	componetHander := v1handler.NewComponentHandler(db.DB)

	v1.GET("/components", componetHander.GetComponent)
	v1.GET("/components/:id", componetHander.GetComponentByID)
	v1.POST("/components/create", componetHander.CreateComponent)

	// Đăng ký các route
	v1.GET("/users", userHandler.GetUser)
	v1.GET("/users/:user_id", userHandler.GetUserById)
	v1.POST("/users/create", userHandler.CreateUser)
	v1.PUT("/users/:user_id", userHandler.PutUser)
	v1.DELETE("/users/:user_id", userHandler.DeleteUser)
	r.Run(":8080")
}
