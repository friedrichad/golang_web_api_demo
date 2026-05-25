package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type ISystemLogRepository interface {
	Create(log *model.SystemLog) error
}

type SystemLogRepository struct {
	DB *gorm.DB
}

var systemLogRepository *SystemLogRepository

func NewSystemLogRepository() *SystemLogRepository {
	if systemLogRepository == nil {
		systemLogRepository = &SystemLogRepository{DB: db.Instance}
	}
	return systemLogRepository
}

func (r *SystemLogRepository) Create(
	log *model.SystemLog,
) error {

	return r.DB.Create(log).Error
}
