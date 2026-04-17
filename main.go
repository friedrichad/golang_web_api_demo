package main

import (
	v1handler "github.com/friedrichad/golang_web_api_demo/internal/api/v1/handler"
	models "github.com/friedrichad/golang_web_api_demo/models"
	"github.com/gin-gonic/gin"
)

var users []models.User

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	userHandler := v1handler.NewUserHandler()
	v1.GET("/users", userHandler.GetUser)
	v1.GET("/users/:user_id", userHandler.GetUserById)
	v1.GET("/users/slug/:slug", userHandler.GetUserSlug)
	v1.POST("/users", userHandler.PostUser)
	v1.PUT("/users/:user_id", userHandler.PutUser)
	v1.DELETE("/users/:user_id", userHandler.DeleteUser)
	r.Run(":8080")

}
