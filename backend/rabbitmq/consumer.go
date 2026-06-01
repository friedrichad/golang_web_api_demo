package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/friedrichad/golang_web_api_demo/backend/model"
	"github.com/friedrichad/golang_web_api_demo/backend/rabbitmq/event"
)

type ISystemLogService interface {
	SaveLog(log model.SystemLogMessage) error
}

func (r *RabbitMQ) ConsumeSystemLogs(svc ISystemLogService) error {
	msgs, err := r.Channel.Consume(
		r.SystemLogQueue.Name,

		"",
		false, // ❗ autoAck = false
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {

			var message model.SystemLogMessage

			err := json.Unmarshal(msg.Body, &message)
			if err != nil {
				log.Println("unmarshal error:", err)
				msg.Nack(false, false)
				continue
			}

			err = svc.SaveLog(message)
			if err != nil {
				log.Println("save db error:", err)
				msg.Nack(false, true) // retry
				continue
			}

			msg.Ack(false)

			log.Println("[SYSTEM LOG SAVED]")
		}
	}()

	log.Println("System log consumer started")

	return nil
}

type IMailService interface {
	SendWelcomeEmail(userID int, email string) error
}

func (r *RabbitMQ) ConsumeMail(mailService IMailService) error {

	msgs, err := r.Channel.Consume(
		MailQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {

			var event event.UserCreatedEvent

			// decode JSON
			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				log.Println("Invalid message:", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf(
				"Received user created event: %+v",
				event,
			)

			log.Printf(
				"Processing mail for user: %d, email: %s\n",
				event.UserID,
				event.Email,
			)

			// 👉 GỌI MAIL SERVICE Ở ĐÂY
			err = mailService.SendWelcomeEmail(
				event.UserID,
				event.Email,
			)

			if err != nil {
				log.Println("Send mail failed:", err)
				msg.Nack(false, true) // retry lại
				continue
			}

			msg.Ack(false)
			log.Println("[WELCOME EMAIL SENT]")
		}
	}()

	log.Println("Mail consumer started")
	return nil
}
