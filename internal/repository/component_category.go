package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type IComponentCategory interface {
	IBaseRepository[model.Request, int]
	GetByComponentCatgeoryById(componentId int) (*model.ComponentCategory, error)
	GetAllByCondition(query model.Request) ([]model.Request, int, error)
	Delete(ids []int) error
	Save(request *model.Request) error
	Update(request *model.Request) error
}