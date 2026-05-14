package cron

import (
	"log"

	"github.com/robfig/cron/v3"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
)

type RequestCron struct {
	service service.IRequestService
}

func NewRequestCron(s service.IRequestService) *RequestCron {
	return &RequestCron{
		service: s,
	}
}

func (c *RequestCron) Start() {
	cr := cron.New()

	// chạy mỗi 1 phút
	cr.AddFunc("@every 1m", func() {
		log.Println("Running ExpireRequests cron...")

		err := c.service.ExpireRequests()
		if err != nil {
			log.Println("ExpireRequests error:", err)
		}
	})

	cr.Start()
}