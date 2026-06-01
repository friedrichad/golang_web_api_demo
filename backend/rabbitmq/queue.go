package rabbitmq

const(
	MailQueue = "mail.queue"
	AuditQueue = "audit.queue"
)

func (r *RabbitMQ) DeclareQueues() error {
	_, err := r.Channel.QueueDeclare(
		MailQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func (r *RabbitMQ) BindQueues() error {
	return r.Channel.QueueBind(
		MailQueue,
		"user.created",
		UserExchange,
		false,
		nil,
	)
}