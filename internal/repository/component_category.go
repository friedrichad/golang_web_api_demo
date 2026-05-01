package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IComponentCategoryRepository interface {
	IBaseRepository[model.ComponentCategory, int32]
	GetAllByCondition(query dtos.ComponentCategoryFilter) ([]model.ComponentCategory, int, error)
	Delete(ids []int32) error
	GetById(id int32) (*model.ComponentCategory, error)
	Save(category *model.ComponentCategory) error
	Update(category *model.ComponentCategory) error
}

type ComponentCategoryRepository struct {
	BaseRepository[model.ComponentCategory, int32]
	DB *gorm.DB
}

var componentCategoryRepository IComponentCategoryRepository

func NewComponentCategoryRepository() IComponentCategoryRepository {
	if componentCategoryRepository == nil {
		componentCategoryRepository = &ComponentCategoryRepository{DB: db.Instance}
		componentCategoryRepository.SetInstance(db.Instance)
	}
	return componentCategoryRepository
}

func (r *ComponentCategoryRepository) GetAllByCondition(query dtos.ComponentCategoryFilter) ([]model.ComponentCategory, int, error) {
	var categories []model.ComponentCategory
	var total int64
	q := r.DB.Model(&model.ComponentCategory{})

	if query.CategoryName != "" {
		q = q.Where("category_name LIKE ?", "%"+query.CategoryName+"%")
	}

	err := q.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Size
	err = q.Offset(offset).Limit(query.Size).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}
	return categories, int(total), nil
}

func (r *ComponentCategoryRepository) Delete(ids []int32) error {
	return r.DB.Where("category_id IN ?", ids).Delete(&model.ComponentCategory{}).Error
}

func (r *ComponentCategoryRepository) GetById(id int32) (*model.ComponentCategory, error) {
	var category *model.ComponentCategory
	err := r.DB.Where("category_id = ?", id).First(&category).Error
	return category, err
}

func (r *ComponentCategoryRepository) Save(category *model.ComponentCategory) error {
	return r.BaseRepository.Create(category)
}

func (r *ComponentCategoryRepository) Update(category *model.ComponentCategory) error {
	return r.BaseRepository.Update(category)
}
