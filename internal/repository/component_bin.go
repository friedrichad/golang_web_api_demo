package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
)

type IComponentBin interface {
	IBaseRepository[model.ComponentBin, int]
	GetByComponentBinId(componentBinId int) (*model.ComponentBin, error)
	GetAuthorities(componentBinId int) ([]string, error)
	GetAllByCondition(query dtos.ComponentBin) ([]model.ComponentBin, int, error)
	Delete(ids []int) error
	Save(componentBin *model.ComponentBin) error
	Update(request *model.Request) error
}