package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IComponentRepository interface {
	IBaseRepository[model.Component, int]
	GetByComponentId(componentId int) (*model.Component, error)
	GetAllByCondition(query dtos.ComponentRequest) ([]model.Component, int, error)
	Delete(ids []int) error
	Save(component *model.Component) error
	Update(component *model.Component) error
}

type ComponentRepository struct{
	BaseRepository[model.Component, int]
	DB *gorm.DB
}

var componentRepository IComponentRepository

func NewComponentRepository() IComponentRepository {
	if  componentRepository == nil {
		componentRepository = &ComponentRepository{DB: db.Instance}
		componentRepository.SetInstance(db.Instance)
	}
	return componentRepository
}

func (c * ComponentRepository) GetByComponentId(componentId int) (*model.Component, error) {
	var component *model.Component
	err := c.DB.Where("component_id = ?", componentId).First(&component).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return component, err
}
func (c *ComponentRepository) GetAllByCondition(query dtos.ComponentRequest) ([]model.Component, int, error) {
	return c.GetPage("Select c.* from component as c "+
	"where (? is Null or c.component_id = ?))"+
	"and (? is Null or c.component_name = ?)) "+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.ComponentID, query.ComponentName, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (c *ComponentRepository) Delete(ids []int) error {
	return c.DB.Exec("delete from component where component_id in ?", ids).Error
}
func (c *ComponentRepository) Save(component *model.Component) error {
	return c.BaseRepository.Create(component)
}
func (c *ComponentRepository) Update(component *model.Component) error {
	return c.BaseRepository.Update(component)
}