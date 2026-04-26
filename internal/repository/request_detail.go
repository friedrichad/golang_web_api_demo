package repository

import (
	// "github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	// "gorm.io/gorm"
)

type IRequestDetailRepository interface {
	IBaseRepository[model.RequestDetail, int]
	GetByRequestId(requestId string) (*model.RequestDetail, error)
	GetAuthorities(requestId int) ([]string, error)
	GetAllByCondition(query model.RequestDetail) ([]model.RequestDetail, int, error)
	Delete(ids []int) error
	Save(request *model.RequestDetail) error
	Update(request *model.RequestDetail) error
}