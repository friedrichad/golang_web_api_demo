package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IBinRepository interface {
	IBaseRepository[model.Bin, int]
	GetAllByCondition(query model.BinFilter) ([]model.Bin, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.Bin, error)
	Save(bin *model.Bin) error
	Update(bin *model.Bin) error
	CreateBinTx(bin *model.Bin) (*model.Bin, error)
	UpdateBinTx(bin *model.Bin) error
	DeleteBinTx(ids []int) error
}

type BinRepository struct {
	BaseRepository[model.Bin, int]
	DB *gorm.DB
}

var binRepository IBinRepository

func NewBinRepository() IBinRepository {
	if binRepository == nil {
		binRepository = &BinRepository{DB: db.Instance}
		binRepository.SetInstance(db.Instance)
	}
	return binRepository
}

func (b *BinRepository) GetAllByCondition(query model.BinFilter) ([]model.Bin, int, error) {
	return b.GetPage("select b.* from bin b"+
		" where (? is null or b.location_in_warehouse like ?)"+
		" and (? is null or b.status_int = ?)"+
		" and (? is null or b.warehouse_id = ?)"+
		" and (? is null or b.created_at >= ?)"+
		" and (? is null or b.created_at <= ?)", query.Page, query.Size, query.LocationInWarehouse, query.LocationInWarehouse, query.StatusInt, query.StatusInt, query.WarehouseID, query.WarehouseID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (b *BinRepository) Delete(ids []int) error {
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

func (b *BinRepository) WithTx(tx *gorm.DB) *BinRepository {
	return &BinRepository{
		BaseRepository: BaseRepository[model.Bin, int]{Instance: tx},
		DB:             tx,
	}
}

// CreateBinTx handles transaction for bin creation
func (b *BinRepository) CreateBinTx(bin *model.Bin) (*model.Bin, error) {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	binRepoTx := b.WithTx(tx)
	err := binRepoTx.Save(bin)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return bin, nil
}

// UpdateBinTx handles transaction for bin update
func (b *BinRepository) UpdateBinTx(bin *model.Bin) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	binRepoTx := b.WithTx(tx)
	err := binRepoTx.Update(bin)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteBinTx handles transaction for bin deletion
func (b *BinRepository) DeleteBinTx(ids []int) error {
	tx := db.Instance.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	binRepoTx := b.WithTx(tx)
	err := binRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
