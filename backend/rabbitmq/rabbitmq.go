package rabbitmq

import (
	"log"

	"github.com/friedrichad/golang_web_api_demo/backend/model/constants"
	"github.com/spf13/viper"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel

	SystemLogQueue amqp.Queue
}

func NewRabbitMQ() (*RabbitMQ, error) {
	url := viper.GetString("rabbitmq.url")

	conn, err := amqp.Dial(url)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = channel.ExchangeDeclare(
		constants.InventoryExchange,

		"topic",

		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	systemLogQueue, err := channel.QueueDeclare(
		constants.SystemLogQueue,

		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = channel.QueueBind(
		systemLogQueue.Name,

		constants.SystemLogRoutingKey,

		constants.InventoryExchange,

		false,
		nil,
	)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Print("RabbitMQ connected")

	return &RabbitMQ{
		Conn:           conn,
		Channel:        channel,
		SystemLogQueue: systemLogQueue,
	}, nil
}
