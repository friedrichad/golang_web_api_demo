package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/friedrichad/golang_web_api_demo/internal/model"
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
