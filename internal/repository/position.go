package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IPositionRepository interface {
	IBaseRepository[model.Position, int]
	GetByPositionName(positionName string) (*model.Position, error)
	GetAll() ([]model.Position, error)
	GetAllByCondition(query model.PositionFilter) ([]model.Position, int, error)
	Delete(ids []int) error
	GetPositionById(positionId int) (*model.Position, error)
	Save(position *model.Position) error
	Update(position *model.Position) error
}

type PositionRepository struct {
	BaseRepository[model.Position, int]
	DB *gorm.DB
}

var positionRepository IPositionRepository

func NewPositionRepository() IPositionRepository {
	if positionRepository == nil {
		positionRepository = &PositionRepository{DB: db.Instance}
		positionRepository.SetInstance(db.Instance)
	}
	return positionRepository
}

func (r *PositionRepository) GetByPositionName(positionName string) (*model.Position, error) {
	var position model.Position
	err := r.DB.Where("position_name = ?", positionName).First(&position).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &position, err
}

func (r *PositionRepository) GetAll() ([]model.Position, error) {
	var positions []model.Position
	err := r.DB.Order("created_at").Find(&positions).Error
	return positions, err
}

func (r *PositionRepository) GetAllByCondition(query model.PositionFilter) ([]model.Position, int, error) {
	return r.GetPage(
		`select p.* 
	 from position p 
	 where (? is null or p.position_name like ?)
	 and (? is null or p.description like ?)
	 and (? is null OR p.created_at >= ?)
	 and (? is null OR p.created_at < ?)`,
		query.Page,
		query.Size,
		query.PositionName,
		query.PositionName,
		query.Description,
		query.Description,
		query.GetDateTo(),
		query.GetDateTo(),
		query.GetDateFrom(),
		query.GetDateFrom(),
	)
}

func (r *PositionRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from position where position_id in ?", ids).Error
}

func (r *PositionRepository) GetPositionById(positionId int) (*model.Position, error) {
	var position model.Position
	err := r.DB.Where("position_id = ?", positionId).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *PositionRepository) Save(position *model.Position) error {
	return r.BaseRepository.Create(position)
}

func (r *PositionRepository) Update(position *model.Position) error {
	return r.BaseRepository.Update(position)
}

func (r *PositionRepository) WithTx(tx *gorm.DB) *PositionRepository {
	return &PositionRepository{
		BaseRepository: BaseRepository[model.Position, int]{Instance: tx},
		DB:             tx,
	}
}
