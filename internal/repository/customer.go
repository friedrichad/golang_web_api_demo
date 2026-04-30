package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type ICustomer interface {
	IBaseRepository[model.Customer, int]
	GetByCustomerId(customerId int) (*model.Customer, error)
	GetAllByCondition(query dtos.CustomerFilter) ([]model.Customer, int, error)
	Delete(ids []int) error
	Save(request *model.Customer) error
	Update(request *model.Customer) error
}
type CustomerRepository struct {
	BaseRepository[model.Customer, int]
	DB *gorm.DB
}

var customerRepository ICustomer

func NewCustomerRepository() ICustomer {
	if customerRepository == nil {
		customerRepository = &CustomerRepository{DB: db.Instance}
		customerRepository.SetInstance(db.Instance)
	}
	return customerRepository
}
func (r *CustomerRepository) GetByCustomerId(customerId int) (*model.Customer, error) {
	var customer *model.Customer
	err := r.DB.Where("customer_id = ?", customerId).First(&customer).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return customer, err
}

func (r *CustomerRepository) GetAllByCondition(query dtos.CustomerFilter) ([]model.Customer, int, error) {
	return r.GetPage("Select c.* from customer as c "+
		"where (? is Null or c.customer_name = ?))"+
		"and (? is null or create_at >= ?) "+
		"and (? is null or create_at < ?) ", query.Page, query.Size, query.CustomerName, query.CustomerName, query.GetDateFrom(), query.GetDateFrom(), query.GetDateTo(), query.GetDateTo())
}
func (r *CustomerRepository) Delete(ids []int) error {
	return r.DB.Exec("delete from customer where customer_id in ?", ids).Error
}
func (r *CustomerRepository) Save(request *model.Customer) error {
	return r.BaseRepository.Create(request)
}
func (r *CustomerRepository) Update(request *model.Customer) error {
	return r.BaseRepository.Update(request)
}
