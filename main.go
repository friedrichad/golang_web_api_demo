package main

import (
	"log"
	"os"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/router"
	"github.com/spf13/viper"
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
	log.SetOutput(os.Stdout)
	redis.InitRedis() 

	router.InitRouter().Run(":" + viper.GetString("port"))
}
