package rabbitmq

const (
	UserExchange      = "user.exchange"
	SystemLogExchange = "system_log.exchange"
)

func (r *RabbitMQ) DeclareExchange() error {
	// Declare UserExchange
	err := r.Channel.ExchangeDeclare(
		UserExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare SystemLogExchange
	return r.Channel.ExchangeDeclare(
		SystemLogExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}
