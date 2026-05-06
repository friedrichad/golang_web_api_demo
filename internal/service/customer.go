package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
)

type ICustomerService interface {
	GetAllCustomers(c *gin.Context) ([]dtos.CustomerResponse, int, *common.Error)
	GetCustomerById(c *gin.Context) (*dtos.CustomerResponse, *common.Error)
	CreateCustomer(c *gin.Context) (*dtos.CustomerResponse, *common.Error)
	UpdateCustomer(c *gin.Context) *common.Error
	DeleteCustomer(c *gin.Context) *common.Error
}

type CustomerService struct {
	customerRepo repository.ICustomer
}

var customerService ICustomerService

func NewCustomerService() ICustomerService {
	if customerService == nil {
		customerService = &CustomerService{
			customerRepo: repository.NewCustomerRepository(),
		}
	}
	return customerService
}

func (s *CustomerService) GetAllCustomers(c *gin.Context) ([]dtos.CustomerResponse, int, *common.Error) {
	var query dtos.CustomerFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	customers, total, err := s.customerRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	res := make([]dtos.CustomerResponse, len(customers))
	for i, cust := range customers {
		res[i] = modelToCustomerResponse(&cust)
	}

	return res, total, nil
}

func (s *CustomerService) GetCustomerById(c *gin.Context) (*dtos.CustomerResponse, *common.Error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	cust, err := s.customerRepo.GetByCustomerId(id)
	if err != nil || cust == nil {
		return nil, common.NotFound
	}

	res := modelToCustomerResponse(cust)
	return &res, nil
}

func (s *CustomerService) CreateCustomer(c *gin.Context) (*dtos.CustomerResponse, *common.Error) {
	var req dtos.CustomerCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return nil, common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	customerRepoTx := s.customerRepo.(*repository.CustomerRepository).WithTx(tx)

	cust := &model.Customer{
		CustomerName: req.CustomerName,
		Phone:        req.Phone,
		Email:        req.Email,
		Address:      req.Address,
		StatusInt:    1,
		CreatedAt:    time.Now(),
	}

	err := customerRepoTx.Save(cust)
	if err != nil {
		tx.Rollback()
		return nil, common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, common.SystemError
	}

	res := modelToCustomerResponse(cust)
	return &res, nil
}

func (s *CustomerService) UpdateCustomer(c *gin.Context) *common.Error {
	var req dtos.CustomerUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	customerRepoTx := s.customerRepo.(*repository.CustomerRepository).WithTx(tx)

	cust, err := customerRepoTx.GetByCustomerId(req.CustomerID)
	if err != nil || cust == nil {
		tx.Rollback()
		return common.NotFound
	}

	if req.CustomerName != "" {
		cust.CustomerName = req.CustomerName
	}
	if req.Phone != "" {
		cust.Phone = req.Phone
	}
	if req.Email != "" {
		cust.Email = req.Email
	}
	if req.Address != "" {
		cust.Address = req.Address
	}
	if req.StatusInt != nil {
		cust.StatusInt = int(*req.StatusInt)
	}
	if req.UpdatedBy != 0 {
		cust.UpdatedBy = int(req.UpdatedBy)
	}
	cust.UpdatedAt = time.Now()

	err = customerRepoTx.Update(cust)
	if err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func (s *CustomerService) DeleteCustomer(c *gin.Context) *common.Error {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	tx := db.Instance.Begin()
	if tx.Error != nil {
		return common.SystemError
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	customerRepoTx := s.customerRepo.(*repository.CustomerRepository).WithTx(tx)

	err := customerRepoTx.Delete(ids)
	if err != nil {
		tx.Rollback()
		return common.SystemError
	}

	if err := tx.Commit().Error; err != nil {
		return common.SystemError
	}

	return nil
}

func modelToCustomerResponse(c *model.Customer) dtos.CustomerResponse {
	return dtos.CustomerResponse{
		CustomerID:   int(c.CustomerID),
		CustomerName: c.CustomerName,
		Phone:        c.Phone,
		Email:        c.Email,
		Address:      c.Address,
		StatusInt:    int(c.StatusInt),
		CreatedBy:    int(c.CreatedBy),
		CreatedAt:    c.CreatedAt,
		UpdatedBy:    int(c.UpdatedBy),
		UpdatedAt:    c.UpdatedAt,
	}
}
