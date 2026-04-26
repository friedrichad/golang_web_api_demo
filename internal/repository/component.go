package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"

)

type IComponentRepository interface {
	IBaseRepository[model.Component, int]
	GetByComponentId(componentId string) (*model.Component, error)
	GetAllByCondition(query model.Component) ([]model.Component, int, error)
	Delete(ids []int) error
	Save(component *model.Component) error
	Update(component *model.Component) error
}

type ComponentRepository struct{
	BaseRepository[model.Component, int]
	DB *gorm.DB
}

var componentRepository ComponentRepository

func NewComponentRepository() IComponentRepository {
	if  componentRepository == nil {
		componentRepository = &ComponentRepository{DB: db.Instance}
		componentRepository.SetInstance(db.Instance)
	}
	return componentRepository
}

func (c * ComponentRepository) GetByComponentId(componentId int) (*model.Component, error) {
	var component model.Component
	err := c.DB.Where("component_id = ?", componentId).First(&component).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return component, err
}
func (c *ComponentRepository) GetAllByCondition(query model.Component) ([]model.Component, int, error) {

}
