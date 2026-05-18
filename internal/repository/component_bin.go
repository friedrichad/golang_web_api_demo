package repository

import (
	"fmt"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IComponentBin interface {
	IBaseRepository[model.ComponentBin, int]
	GetByComponentBinId(componentBinId int) (*model.ComponentBin, error)
	GetByComponentAndBinId(componentID int, binID int) (*model.ComponentBin, error)
	GetByComponentAndBinIds(componentIDs []int, binIDs []int) map[string]*model.ComponentBin
	GetAllByCondition(query model.ComponentBinFilter) ([]model.ComponentBin, int, error)
	Delete(ids []int) error
	Save(componentBin *model.ComponentBin) error
	Update(request *model.ComponentBin) error
	CreateBatch(bins []model.ComponentBin) error
	UpdateBatch(bins []model.ComponentBin) error
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

func (r *ComponentBinRepository) GetByComponentBinId(componentBinId int) (*model.ComponentBin, error) {
	var componentBin *model.ComponentBin
	err := r.DB.Where("component_bin_id = ?", componentBinId).First(&componentBin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return componentBin, err
}

func (r *ComponentBinRepository) GetByComponentAndBinId(componentID int, binID int) (*model.ComponentBin, error) {
	var componentBin *model.ComponentBin
	err := r.DB.Where("component_id = ? AND bin_id = ?", componentID, binID).First(&componentBin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return componentBin, err
}

func (r *ComponentBinRepository) GetAllByCondition(query model.ComponentBinFilter) ([]model.ComponentBin, int, error) {
	return r.GetPage("Select cb.* from component_bin as cb "+
		"where (? is Null or cb.quantity = ?))"+
		"and (? is Null or cb.component_id = ?)) "+
		"and (? is Null or cb.bin_id = ?)) "+
		"and (? is Null or cb.created_at >= ?)"+
		"and (? is null or cb.created_at < ?) ", query.Page, query.Size, query.Quantity, query.Quantity, query.ComponentID, query.ComponentID, query.BinID, query.BinID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (r *ComponentBinRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from component_bin where component_bin_id in ?", ids).Error
}

func (r *ComponentBinRepository) Save(componentBin *model.ComponentBin) error {
	return r.BaseRepository.Create(componentBin)
}

func (r *ComponentBinRepository) Update(request *model.ComponentBin) error {
	return r.BaseRepository.Update(request)
}

// GetByComponentAndBinIds loads all component bins matching component and bin IDs in one query
func (r *ComponentBinRepository) GetByComponentAndBinIds(componentIDs []int, binIDs []int) map[string]*model.ComponentBin {
	var bins []model.ComponentBin
	r.DB.Where("component_id IN ? AND bin_id IN ?", componentIDs, binIDs).Find(&bins)

	// Map by "component_id_bin_id" for easy lookup
	result := make(map[string]*model.ComponentBin)
	for i := range bins {
		key := formatBinKey(bins[i].ComponentID, bins[i].BinID)
		result[key] = &bins[i]
	}
	return result
}

// CreateBatch creates multiple component bins in a single query
func (r *ComponentBinRepository) CreateBatch(bins []model.ComponentBin) error {
	if len(bins) == 0 {
		return nil
	}
	return r.DB.CreateInBatches(bins, 100).Error
}

// UpdateBatch updates multiple component bins efficiently
func (r *ComponentBinRepository) UpdateBatch(bins []model.ComponentBin) error {
	if len(bins) == 0 {
		return nil
	}
	// Update in batches for efficiency
	for i := 0; i < len(bins); i += 100 {
		end := i + 100
		if end > len(bins) {
			end = len(bins)
		}
		// Use GORM's Save which performs UPDATE
		if err := r.DB.Save(bins[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}

// formatBinKey creates a key for component-bin mapping
func formatBinKey(componentID int, binID int) string {
	return fmt.Sprintf("%d_%d", componentID, binID)
}

func (r *ComponentBinRepository) WithTx(tx *gorm.DB) *ComponentBinRepository {
	return &ComponentBinRepository{
		BaseRepository: BaseRepository[model.ComponentBin, int]{Instance: tx},
		DB:             tx,
	}
}
