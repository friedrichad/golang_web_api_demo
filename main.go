package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/demo", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, this is a json response!",
		})
	})
	r.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"users": "List of users",
		})
	})
	r.GET("/users/:user_id", func(ctx *gin.Context) {
		user_id := ctx.Param("user_id")
		ctx.JSON(200, gin.H{
			"data":    "User details",
			"user_id": user_id,
		})
	})
	r.GET("/product/:product_name", func(ctx *gin.Context) {
		product_name := ctx.Param("product_name")
		price := ctx.Query("price")
		ctx.JSON(200, gin.H{
			"data":         "Product details",
			"product_name": product_name,
			"price":        price,
		})
	})
	r.Run(":8080")
}
