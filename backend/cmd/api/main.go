package main

import (
	"log"
	"os"

	"github.com/friedrichad/golang_web_api_demo/backend/configs/db"
	"github.com/friedrichad/golang_web_api_demo/backend/configs/mail"
	"github.com/friedrichad/golang_web_api_demo/backend/cron"
	"github.com/friedrichad/golang_web_api_demo/backend/rabbitmq"
	"github.com/friedrichad/golang_web_api_demo/backend/redis"
	"github.com/friedrichad/golang_web_api_demo/backend/router"
	"github.com/friedrichad/golang_web_api_demo/backend/service"
	"github.com/spf13/viper"
)

// @title Auth Server API
// @version 1.0
// @description JWT Auth Server
// @host localhost:8080
// @BasePath /

func main() {
	// Try to find config file in common locations
	configPaths := []string{
		"configs/config.yaml",
		"../../configs/config.yaml",
		"../configs/config.yaml",
	}

	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		log.Fatal("config.yaml not found in any expected location")
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	db.InitMysql()
	log.SetOutput(os.Stdout)
	redis.InitRedis()
	mail.InitMailConfig()
	var rmq *rabbitmq.RabbitMQ

	rmqInstance, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		log.Println(
			"RabbitMQ unavailable:",
			err,
		)
	} else {

		rmq = rmqInstance

		err = rmq.DeclareExchange()
		if err != nil {
			log.Fatal(err)
		}

		err = rmq.DeclareQueues()
		if err != nil {
			log.Fatal(err)
		}

		err = rmq.BindQueues()
		if err != nil {
			log.Fatal(err)
		}
	}
	requestService := service.NewRequestService()
	systemLogService := service.NewSystemLogService()
	if rmq != nil {
		go rmq.ConsumeSystemLogs(systemLogService)
		mailService := service.NewMailService()
		go rmq.ConsumeMail(mailService)
	}
	requestCron := cron.NewRequestCron(
		requestService,
	)
	requestCron.Start()
	router.InitRouter(rmq).Run(
		":" + viper.GetString("port"),
	)
}
