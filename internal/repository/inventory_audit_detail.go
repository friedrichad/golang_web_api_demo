package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type IInventoryAuditDetail interface {
	IBaseRepository[model.Request, int]
	GetByRequestId(requestId string) (*model.Request, error)
	GetAuthorities(requestId int) ([]string, error)
	GetAllByCondition(query model.Request) ([]model.Request, int, error)
	Delete(ids []int) error
	Save(request *model.Request) error
	Update(request *model.Request) error
}