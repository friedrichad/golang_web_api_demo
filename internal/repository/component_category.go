package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IComponentCategoryRepository interface {
	IBaseRepository[model.ComponentCategory, int]
	GetAllByCondition(query model.ComponentCategoryFilter) ([]model.ComponentCategory, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.ComponentCategory, error)
	Save(category *model.ComponentCategory) error
	Update(category *model.ComponentCategory) error
	CreateComponentCategoryTx(category *model.ComponentCategory) (*model.ComponentCategory, error)
	UpdateComponentCategoryTx(category *model.ComponentCategory) error
	DeleteComponentCategoryTx(ids []int) error
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

func (r *ComponentCategoryRepository) GetAllByCondition(query model.ComponentCategoryFilter) ([]model.ComponentCategory, int, error) {
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

// CreateComponentCategoryTx handles transaction for component category creation
func (r *ComponentCategoryRepository) CreateComponentCategoryTx(category *model.ComponentCategory) (*model.ComponentCategory, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	categoryRepoTx := r.WithTx(tx)
	err := categoryRepoTx.Save(category)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateComponentCategoryTx handles transaction for component category update
func (r *ComponentCategoryRepository) UpdateComponentCategoryTx(category *model.ComponentCategory) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	categoryRepoTx := r.WithTx(tx)
	err := categoryRepoTx.Update(category)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteComponentCategoryTx handles transaction for component category deletion
func (r *ComponentCategoryRepository) DeleteComponentCategoryTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	categoryRepoTx := r.WithTx(tx)
	err := categoryRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
