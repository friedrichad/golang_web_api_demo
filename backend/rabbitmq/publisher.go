package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/friedrichad/golang_web_api_demo/backend/model/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Publish(
	routingKey string,
	data interface{},
) error {

	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return r.Channel.PublishWithContext(
		context.Background(),

		constants.InventoryExchange,

		routingKey,

		false,
		false,

		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
