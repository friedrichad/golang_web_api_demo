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
	GetAllByCondition(query dtos.ComponentFilter) ([]model.Component, int, error)
	Delete(ids []int) error
	Save(component *model.Component) error
	Update(component *model.Component) error
	CreateComponentTx(component *model.Component, categories []model.ComponentCategoryMap) (*model.Component, error)
	UpdateComponentTx(component *model.Component, categories []model.ComponentCategoryMap) error
	DeleteComponentTx(ids []int) error
}

type ComponentRepository struct {
	BaseRepository[model.Component, int]
	DB *gorm.DB
}

var componentRepository IComponentRepository

func NewComponentRepository() IComponentRepository {
	if componentRepository == nil {
		componentRepository = &ComponentRepository{DB: db.Instance}
		componentRepository.SetInstance(db.Instance)
	}
	return componentRepository
}

func (c *ComponentRepository) GetByComponentId(componentId int) (*model.Component, error) {
	var component *model.Component
	err := c.DB.Where("component_id = ?", componentId).First(&component).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return component, err
}
func (c *ComponentRepository) GetAllByCondition(query dtos.ComponentFilter) ([]model.Component, int, error) {
	return c.GetPage("Select c.* from component c "+
		" where (? is Null or c.component_id = ?)"+
		" and (? is Null or c.component_name = ?) "+
		" and (? is null or created_at >= ?) "+
		" and (? is null or created_at <= ?) ", query.Page, query.Size, query.ComponentID, query.ComponentID, query.ComponentName, query.ComponentName, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
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

func (c *ComponentRepository) WithTx(tx *gorm.DB) *ComponentRepository {
	return &ComponentRepository{
		BaseRepository: BaseRepository[model.Component, int]{Instance: tx},
		DB:             tx,
	}
}

// CreateComponentTx handles transaction for component creation with category mappings
func (c *ComponentRepository) CreateComponentTx(component *model.Component, categories []model.ComponentCategoryMap) (*model.Component, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	componentRepoTx := c.WithTx(tx)
	err := componentRepoTx.Save(component)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, cat := range categories {
		cat.ComponentID = component.ComponentID
		if err := tx.Create(&cat).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return component, nil
}

// UpdateComponentTx handles transaction for component update with category mappings
func (c *ComponentRepository) UpdateComponentTx(component *model.Component, categories []model.ComponentCategoryMap) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	componentRepoTx := c.WithTx(tx)
	err := componentRepoTx.Update(component)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing category mappings
	if err := tx.Where("component_id = ?", component.ComponentID).Delete(&model.ComponentCategoryMap{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create new category mappings
	for _, cat := range categories {
		cat.ComponentID = component.ComponentID
		if err := tx.Create(&cat).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteComponentTx handles transaction for component deletion with category mappings
func (c *ComponentRepository) DeleteComponentTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Delete category mappings first
	if err := tx.Where("component_id IN ?", ids).Delete(&model.ComponentCategoryMap{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete components
	componentRepoTx := c.WithTx(tx)
	err := componentRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
