package main

import (
	config "github.com/friedrichad/golang_web_api_demo/configs"
	db "github.com/friedrichad/golang_web_api_demo/db"
	v1handler "github.com/friedrichad/golang_web_api_demo/internal/api/v1/handler"
	"github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.LoadConfig("configs/config.json")
    if err != nil {
        panic(err)
    }

    db.ConnectDB(cfg)

    r := gin.Default()
    v1 := r.Group("/api/v1")

    // Khởi tạo UserHandler với DB đã connect
    userHandler := v1handler.NewUserHandler(db.DB)

    // Đăng ký các route
    v1.GET("/users", userHandler.GetUser)
    v1.GET("/users/:user_id", userHandler.GetUserById)
    v1.POST("/users", userHandler.PostUser)
    v1.PUT("/users/:user_id", userHandler.PutUser)
    v1.DELETE("/users/:user_id", userHandler.DeleteUser)
	r.Run(":8080")
}
