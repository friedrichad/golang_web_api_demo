package rabbitmq

const(
	UserExchange = "user.exchange"
)

func (r *RabbitMQ) DeclareExchange() error {
	return r.Channel.ExchangeDeclare(
		UserExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}