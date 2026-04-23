package repository

import (
	"gorm.io/gorm"
	"strings"
)

type BaseRepository[E any, T any] struct {
	Instance *gorm.DB
}

type IBaseRepository[E any, T any] interface {
	SetInstance(instance *gorm.DB)
	Delete(ids []T) error
	GetAll() ([]E, error)
	GetById(id T) (*E, error)
	GetPage(sql string, page int, size int, values ...interface{}) ([]E, int, error)
}

func (r *BaseRepository[E, T]) SetInstance(instance *gorm.DB) {
	r.Instance = instance
}

func (r *BaseRepository[E, T]) Delete(ids []T) error {
	var e []E
	return r.Instance.Delete(&e, ids).Error
}

func (r *BaseRepository[E, T]) GetAll() ([]E, error) {
	var slice []E
	err := r.Instance.Find(&slice).Error
	return slice, err
}

func (r *BaseRepository[E, T]) GetById(id T) (*E, error) {
	var e *E
	err := r.Instance.First(&e, "id = ?", id).Error
	return e, err
}

func (r *BaseRepository[E, T]) Create(e *E) error {
	return r.Instance.Create(&e).Error
}

func (r *BaseRepository[E, T]) Update(e *E) error {
	return r.Instance.Model(&e).Updates(e).Error
}

func (r *BaseRepository[E, T]) GetPage(sql string, page int, size int, values ...interface{}) ([]E, int, error) {
	if size == 0 {
		size = 10
	}
	offset := page * size
	var slice []E
	var err error
	if values == nil {
		err = r.Instance.Raw(sql+" limit ? offset ?", size, offset).Scan(&slice).Error
	} else {
		err = r.Instance.Raw(sql+" limit ? offset ?", append(values, size, offset)...).Scan(&slice).Error
	}
	if err != nil {
		return nil, 0, err
	}
	sql = "select count(1) " + sql[strings.LastIndex(sql, "from"):]
	var total int
	if values == nil {
		err = r.Instance.Raw(sql).Scan(&total).Error
	} else {
		err = r.Instance.Raw(sql, values...).Scan(&total).Error
	}
	return slice, total, err
}
