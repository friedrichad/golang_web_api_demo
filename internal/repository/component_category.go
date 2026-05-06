package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IComponentCategoryRepository interface {
	IBaseRepository[model.ComponentCategory, int]
	GetAllByCondition(query dtos.ComponentCategoryFilter) ([]model.ComponentCategory, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.ComponentCategory, error)
	Save(category *model.ComponentCategory) error
	Update(category *model.ComponentCategory) error
}

type ComponentCategoryRepository struct {
	BaseRepository[model.ComponentCategory, int]
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
	return r.GetPage("select cc.* from ComponentCategory cc"+
		" where (? is null or cc.category_name like ?)"+
		" and (? is null or cc.created_at >= ?)"+
		" and (? is null or cc.created_at <= ?)", query.Page, query.Size, query.CategoryName, query.CategoryName, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (r *ComponentCategoryRepository) Delete(ids []int) error {
	return r.DB.Where("category_id in ?", ids).Delete(&model.ComponentCategory{}).Error
}

func (r *ComponentCategoryRepository) GetById(id int) (*model.ComponentCategory, error) {
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

func (r *ComponentCategoryRepository) WithTx(tx *gorm.DB) *ComponentCategoryRepository {
	return &ComponentCategoryRepository{
		BaseRepository: BaseRepository[model.ComponentCategory, int]{Instance: tx},
		DB:             tx,
	}
}
