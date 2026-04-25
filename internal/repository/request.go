package repository

import "github.com/friedrichad/golang_web_api_demo/internal/model"

type IRequestRepository interface {
	IBaseRepository[model.Request, int]
	GetByRequestId(requestId string) (*model.Request, error)
	GetAuthorities(userId int) ([]string, error)
	GetAllByCondition(query model.Request) ([]model.User, int, error)
	Delete(ids []int) error
	GetByUuid(id int) (*model.Request, error)
	Save(user *model.Request) error
	Update(user *model.UserUpdate) error
}