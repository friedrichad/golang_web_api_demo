package repository

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"gorm.io/gorm"
)

type IRequestRepository interface {
	IBaseRepository[model.Request, int]
	GetByRequestId(requestId int) (*model.Request, error)
	GetAllByCondition(query dtos.RequestFilter) ([]model.Request, int, error)
	Delete(ids []int) error
	Save(request *model.Request) error
	Update(request *model.Request) error
	CanApprove(approverId int, requesterId int) (bool, error)
	GetExpiredPendingRequests() ([]model.Request, error)
}

type RequestRepository struct {
	BaseRepository[model.Request, int]
	DB *gorm.DB
}

var requestRepository IRequestRepository

func NewRequestRepository() IRequestRepository {
	if requestRepository == nil {
		requestRepository = &RequestRepository{DB: db.Instance}
		requestRepository.SetInstance(db.Instance)
	}
	return requestRepository
}

func (r *RequestRepository) GetByRequestId(requestId int) (*model.Request, error) {
	var request *model.Request
	err := r.DB.Where("request_id = ?", requestId).First(&request).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return request, err
}
func (r *RequestRepository) GetAllByCondition(query dtos.RequestFilter) ([]model.Request, int, error) {
	return r.GetPage("Select r.* from request as r "+
		" where (? is Null or r.request_id = ?)"+
		" and (? is Null or r.request_type = ?) "+
		" and (?  is Null or r.status_int = ?) "+
		" and (? is null or created_at >= ?) "+
		" and (? is null or created_at < ?) ", query.Page, query.Size, query.RequestID, query.RequestID, query.RequestType, query.RequestType, query.StatusInt, query.StatusInt, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
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

func (r *RequestRepository) WithTx(tx *gorm.DB) *RequestRepository {
	return &RequestRepository{
		BaseRepository: BaseRepository[model.Request, int]{Instance: tx},
		DB:             tx,
	}
}

func (r *RequestRepository) CanApprove(approverId int, requesterId int) (bool, error) {
	var result int
	err := r.DB.Raw("SELECT 1"+
		" from user u_approver"+
		" join user u_requester on u_requester.user_id = ?"+
		" join position_hierarchy ph"+
		" on ph.ancestor_id = u_approver.position_id"+
		" and ph.descendant_id = u_requester.position_id"+
		" where u_approver.user_id = ?"+
		" and ph.depth > 0"+
		" limit 1", requesterId, approverId).Scan(&result).Error
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (r *RequestRepository) GetExpiredPendingRequests() ([]model.Request, error) {
	var requests []model.Request
	err := r.DB.Where("status_int = ? AND expired_date < ?", constants.RequestStatusPending, time.Now()).Find(&requests).Error
	return requests, err
}
