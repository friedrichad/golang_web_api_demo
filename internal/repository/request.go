package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"gorm.io/gorm"
)

type IRequestRepository interface {
	IBaseRepository[model.Request, int]
	GetByRequestId(requestId string) (*model.Request, error)
	GetAllByCondition(query dtos.RequestRequest) ([]model.Request, int, error)
	Delete(ids []int) error
	Save(request *model.Request) error
	Update(request *model.Request) error
}

type RequestRepository struct{
	BaseRepository[model.Request, int]
	DB * gorm.DB
}

var requestRepository IRequestRepository

func NewRequestRepository() IRequestRepository {
	if requestRepository == nil {
		requestRepository = &RequestRepository{DB: db.Instance}
		requestRepository.SetInstance(db.Instance)
	}
	return requestRepository
}

func (r *RequestRepository) GetByRequestId(requestId string) (*model.Request, error) {
	var request *model.Request
	err := r.DB.Where("request_id = ?", requestId).First(&request).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return request, err
}
func (r *RequestRepository) GetAllByCondition(query dtos.RequestRequest) ([]model.Request, int, error) {
	return r.GetPage("Select r.* from request as r "+
	"where (? is Null or r.request_id = ?))"+
	" and (? is Null or r.request_type = ?)) "+ 
	"and (?  is Null or r.status_int = ?)) "+
	"and (? is null or create_time >= ?) "+
	"and (? is null or create_time < ?) ", query.Page, query.Size, query.RequestID, query.RequestID, query.RequestType, query.RequestType, query.StatusInt, query.StatusInt, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (r *RequestRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from request where request_id in ?", ids).Error
}
func (r *RequestRepository) Save(request *model.Request) error {
	return r.BaseRepository.Create(request)
}
func (r *RequestRepository) Update(request *model.Request) error {
	return r.BaseRepository.Update(request)
}