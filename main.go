package main

import (
	"log"
	"os"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/cron"
	"github.com/friedrichad/golang_web_api_demo/internal/rabbitmq"
	"github.com/friedrichad/golang_web_api_demo/internal/redis"
	"github.com/friedrichad/golang_web_api_demo/internal/router"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/spf13/viper"
)

// @title Auth Server API
// @version 1.0
// @description JWT Auth Server
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	viper.SetConfigFile(
		"internal/configs/config.yaml",
	)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	db.InitMysql()
	log.SetOutput(os.Stdout)
	redis.InitRedis()
	var rmq *rabbitmq.RabbitMQ
	rmqInstance, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		log.Println(
			"RabbitMQ unavailable:",
			err,
		)
	} else {
		rmq = rmqInstance
	}
	requestService := service.NewRequestService()
	systemLogService := service.NewSystemLogService()
	if rmq != nil {
		go rmq.ConsumeSystemLogs(systemLogService)
	}
	requestCron := cron.NewRequestCron(
		requestService,
	)
	requestCron.Start()
	router.InitRouter(rmq).Run(
		":" + viper.GetString("port"),
	)
}
