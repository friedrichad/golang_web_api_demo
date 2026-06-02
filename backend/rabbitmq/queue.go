package rabbitmq

const (
	MailQueue      = "mail.queue"
	AuditQueue     = "audit.queue"
	SystemLogQueue = "system_log.queue"
)

func (r *RabbitMQ) DeclareQueues() error {
	// Declare MailQueue
	_, err := r.Channel.QueueDeclare(
		MailQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare SystemLogQueue
	_, err = r.Channel.QueueDeclare(
		SystemLogQueue,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func (r *RabbitMQ) BindQueues() error {
	// Bind MailQueue to UserExchange
	err := r.Channel.QueueBind(
		MailQueue,
		"user.created",
		UserExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind SystemLogQueue to SystemLogExchange
	return r.Channel.QueueBind(
		SystemLogQueue,
		"system.*",
		SystemLogExchange,
		false,
		nil,
	)
}
