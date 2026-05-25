package service

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type SystemLogService struct {
}

func NewSystemLogService() *SystemLogService {
	return &SystemLogService{}
}

func (s *SystemLogService) SaveLog(
	message model.SystemLogMessage,
) error {

	systemLog := model.SystemLog{
		UserID: message.UserID,

		HTTPMethod: message.HTTPMethod,

		Route: message.Route,

		StatusInt: message.StatusInt,

		IPAddress: message.IPAddress,

		ResponseBody: message.ResponseBody,

		ExecutedAt: time.Now(),
	}

	return db.Instance.Create(
		&systemLog,
	).Error
}
