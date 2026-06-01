package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish message to UserExchange (for user events like user.created)
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

		UserExchange,

		routingKey,

		false,
		false,

		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
