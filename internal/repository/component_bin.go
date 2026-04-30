package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"gorm.io/gorm"

)

type IComponentBin interface {
	IBaseRepository[model.ComponentBin, int]
	GetByComponentBinId(componentBinId int) (*model.ComponentBin, error)
	GetAllByCondition(query dtos.ComponetBinFilter) ([]model.ComponentBin, int, error)
	Delete(ids []int) error
	Save(componentBin *model.ComponentBin) error
	Update(request *model.ComponentBin) error
}

type ComponentBinRepository struct {
	BaseRepository[model.ComponentBin, int]
	DB *gorm.DB
}

var componentBinRepository IComponentBin
func NewComponentBinRepository() IComponentBin {
	if componentBinRepository == nil {
		componentBinRepository = &ComponentBinRepository{DB: db.Instance}
		componentBinRepository.SetInstance(db.Instance)
	}
	return componentBinRepository
}

func (r *ComponentBinRepository) GetByComponentBinId(componentBinId int) (*model.ComponentBin, error){
	var componentBin *model.ComponentBin
	err := r.DB.Where("component_bin_id = ?", componentBinId).First(&componentBin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return componentBin, err
}

func (r *ComponentBinRepository) GetAllByCondition(query dtos.ComponetBinFilter) ([]model.ComponentBin, int, error) {
	return r.GetPage("Select cb.* from component_bin as cb "+
		"where (? is Null or cb.quantity = ?))"+
		"and (? is Null or cb.component_id = ?)) "+
		"and (? is Null or cb.bin_id = ?)) "+
		"and (? is Null or cb.created_at >= ?)"+
		"and (? is null or cb.created_at < ?) ", query.Page, query.Size, query.Quantity, query.Quantity, query.ComponentID, query.ComponentID, query.BinID, query.BinID, query.GetDateFrom(), query.GetDateFrom(),query.GetDateTo(), query.GetDateTo())
}

func (r *ComponentBinRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from component_bin where component_bin_id in ?", ids).Error
}

func (r *ComponentBinRepository) Save(componentBin *model.ComponentBin) error{
	return r.BaseRepository.Create(componentBin)
}

func (r *ComponentBinRepository) Update(request *model.ComponentBin) error {
	return r.BaseRepository.Update(request)
}