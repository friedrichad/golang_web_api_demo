package controller

import (
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/gin-gonic/gin"
)

type ICustomerController interface {
	GetAllCustomers() gin.HandlerFunc
	GetCustomerById() gin.HandlerFunc
	CreateCustomer() gin.HandlerFunc
	UpdateCustomer() gin.HandlerFunc
	DeleteCustomer() gin.HandlerFunc
}

type CustomerController struct {
	BaseController[model.CustomerResponse]
	customerService service.ICustomerService
}

func NewCustomerController() ICustomerController {
	customerService := service.NewCustomerService()
	return &CustomerController{customerService: customerService}
}

func (c *CustomerController) GetAllCustomers() gin.HandlerFunc {
	return c.ResponsePage(c.customerService.GetAllCustomers)
}

func (c *CustomerController) GetCustomerById() gin.HandlerFunc {
	return c.ResponsePointer(c.customerService.GetCustomerById)
}

func (c *CustomerController) CreateCustomer() gin.HandlerFunc {
	return c.ResponsePointer(c.customerService.CreateCustomer)
}

func (c *CustomerController) UpdateCustomer() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.customerService.UpdateCustomer)
}

func (c *CustomerController) DeleteCustomer() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.customerService.DeleteCustomer)
}
