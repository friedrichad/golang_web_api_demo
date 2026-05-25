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

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Get a list of all customers with optional filtering and pagination
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security ApiKeyAuth
// @Success 200 {object} []model.CustomerResponse
// @Failure 400 {object} common.Error
// @Router /customers [get]
func (c *CustomerController) GetAllCustomers() gin.HandlerFunc {
	return c.ResponsePage(c.customerService.GetAllCustomers)
}

// GetCustomerById godoc
// @Summary Get a customer by ID
// @Description Retrieve details of a specific customer by their ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.CustomerResponse
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /customers/{id} [get]
func (c *CustomerController) GetCustomerById() gin.HandlerFunc {
	return c.ResponsePointer(c.customerService.GetCustomerById)
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Add a new customer to the system
// @Tags customers
// @Accept json
// @Produce json
// @Param request body model.CustomerCreate true "Customer Info"
// @Security ApiKeyAuth
// @Success 200 {object} model.CustomerResponse
// @Failure 400 {object} common.Error
// @Router /customers [post]
func (c *CustomerController) CreateCustomer() gin.HandlerFunc {
	return c.ResponsePointer(c.customerService.CreateCustomer)
}

// UpdateCustomer godoc
// @Summary Update an existing customer
// @Description Modify details of an existing customer
// @Tags customers
// @Accept json
// @Produce json
// @Param request body model.CustomerUpdate true "Customer Update Info"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Failure 404 {object} common.Error
// @Router /customers [put]
func (c *CustomerController) UpdateCustomer() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.customerService.UpdateCustomer)
}

// DeleteCustomer godoc
// @Summary Delete customers
// @Description Delete one or multiple customers by their IDs
// @Tags customers
// @Accept json
// @Produce json
// @Param request body []int true "Array of Customer IDs to delete"
// @Security ApiKeyAuth
// @Success 200 "Success"
// @Failure 400 {object} common.Error
// @Router /customers [delete]
func (c *CustomerController) DeleteCustomer() gin.HandlerFunc {
	return c.ResponseSuccessOnly(c.customerService.DeleteCustomer)
}
