package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
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

	log.Print("RabbitMQ connected")

	return &RabbitMQ{
		Conn:    conn,
		Channel: channel,
	}, nil
}
