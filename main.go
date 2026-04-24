package main

import (
	"github.com/friedrichad/golang_web_api_demo/internal/router"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/spf13/viper"
	"log"
)

// @title Auth Server API
// @version 1.0
// @description JWT Auth Server
// @host localhost:8080
// @BasePath /

func main() {	
	viper.SetConfigFile("internal/configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	db.InitMysql()

	router.InitRouter().Run(":" + viper.GetString("port"))
}
