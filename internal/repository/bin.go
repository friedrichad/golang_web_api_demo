package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IBinRepository interface {
	IBaseRepository[model.Bin, int]
	GetAllByCondition(query dtos.BinFilter) ([]model.Bin, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.Bin, error)
	Save(bin *model.Bin) error
	Update(bin *model.Bin) error
}

type BinRepository struct {
	BaseRepository[model.Bin, int]
	DB *gorm.DB
}

var binRepository IBinRepository

func NewBinRepository() IBinRepository{
	if binRepository == nil {
		binRepository = &BinRepository{DB: db.Instance}
		binRepository.SetInstance(db.Instance)
	}
	return binRepository
}

func (b *BinRepository) GetAllByCondition(query dtos.BinFilter) ([]model.Bin, int, error) {
	return b.GetPage("SELECT b.*, w.name as warehouse_name FROM bin b LEFT JOIN warehouse w ON b.warehouse_id = w.id WHERE b.location_in_warehouse LIKE ? AND b.status_int = ? AND b.warehouse_id = ?", query.Page, query.Size,"%"+query.LocationInWarehouse+"%", query.StatusInt, query.WarehouseID)
}
func (b* BinRepository) Delete(ids []int) error {
	return b.DB.Exec("delete from bin where bin_id in ?", ids).Error
}
func (b *BinRepository) GetById(id int) (*model.Bin, error) {
	var bin *model.Bin
	err := b.DB.Where("bin_id = ?", id).First(&bin).Error
	return bin, err
}
func (b *BinRepository) Save(bin *model.Bin) error {
	return b.BaseRepository.Create(bin)
}
func (b *BinRepository) Update(bin *model.Bin) error {
	return b.BaseRepository.Update(bin)
}