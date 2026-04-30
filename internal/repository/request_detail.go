package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IRequestDetailRepository interface {
	IBaseRepository[model.RequestDetail, int]
	GetByRequestId(requestId int) (*model.RequestDetail, error)
	GetAllByCondition(query dtos.RequestDetailRequest) ([]model.RequestDetail, int, error)
	Delete(ids []int) error
	Save(request *model.RequestDetail) error
	Update(request *model.RequestDetail) error
}
type RequestDetailRepository struct{
	BaseRepository[model.RequestDetail, int]
	DB *gorm.DB
}

var requestDetailRepository IRequestDetailRepository

func NewRequestDetailRepository() IRequestDetailRepository{
	if requestDetailRepository == nil{
		requestDetailRepository = &RequestDetailRepository{DB: db.Instance}
		requestDetailRepository.SetInstance(db.Instance)	
	}
	return requestDetailRepository
}
func (r *RequestDetailRepository) GetByRequestId(requestId int) (*model.RequestDetail, error) {
	var requestDetail *model.RequestDetail
	err := r.DB.Where("request_id = ?", requestId).First(&requestDetail).Error
	if err != nil {
		return nil, err
	}
	return requestDetail, nil
}
func (r *RequestDetailRepository) GetAllByCondition(query dtos.RequestDetailRequest) ([]model.RequestDetail, int, error) {
	return r.GetPage("Select rd.* from request_detail as rd "+
	"where (? is Null or rd.request_id = ?))"+
	"request_id = ?)) "+
	"and (? is Null or rd.component_id = ?)) "+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.RequestID, query.RequestID, query.ComponentID, query.ComponentID, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}

func (r *RequestDetailRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from request_detail where request_detail_id in ?", ids).Error
}
func (r *RequestDetailRepository) Save(request *model.RequestDetail) error {
	return r.BaseRepository.Create(request)
}
func (r *RequestDetailRepository) Update(request *model.RequestDetail) error {
	return r.BaseRepository.Update(request)
}