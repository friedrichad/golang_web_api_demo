package controller

import (
	"net/http"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/gin-gonic/gin"
)
type BaseController[T any] struct {
}

type IBaseController interface {
}

// Success
/**
* Returns a success response
 */
func (controller *BaseController[T]) Success(c *gin.Context, v any) {
	c.JSON(http.StatusOK, model.ResponseWrapper{
		Code:    common.Success.Code,
		Message: common.Success.Message,
		Data:    v,
	})
}

// Error
/**
* Returns a error response
 */
func (controller *BaseController[T]) Error(c *gin.Context, err *common.Error, v any) {
	c.JSON(http.StatusOK, model.ResponseWrapper{
		Code:    err.Code,
		Message: err.Message,
		Data:    v,
	})
}

// ResponseArray
/**
* Return function response a slice
 */
func (controller *BaseController[T]) ResponseArray(serviceFunc func(g *gin.Context) ([]T, *common.Error)) gin.HandlerFunc {
	return func(g *gin.Context) {
		body, err := serviceFunc(g)
		if err != nil {
			controller.Error(g, err, nil)
			return
		}
		controller.Success(g, body)
	}
}

// ResponsePage
/**
* Return function response a page
 */
func (controller *BaseController[T]) ResponsePage(serviceFunc func(g *gin.Context) ([]T, int, *common.Error)) gin.HandlerFunc {
	return func(g *gin.Context) {
		content, total, err := serviceFunc(g)
		if err != nil {
			controller.Error(g, err, nil)
			return
		}
		controller.Success(g, model.Page[T]{Content: content, Total: total})
	}
}

// ResponseObject
/**
* Return function response a object
 */
func (controller *BaseController[T]) ResponseObject(serviceFunc func(g *gin.Context) (T, *common.Error)) gin.HandlerFunc {
	return func(g *gin.Context) {
		body, err := serviceFunc(g)
		if err != nil {
			controller.Error(g, err, nil)
			return
		}
		controller.Success(g, body)
	}
}

// ResponsePointer
/**
* Return function response a pointer of object
 */
func (controller *BaseController[T]) ResponsePointer(serviceFunc func(g *gin.Context) (*T, *common.Error)) gin.HandlerFunc {
	return func(g *gin.Context) {
		body, err := serviceFunc(g)
		if err != nil {
			controller.Error(g, err, nil)
			return
		}
		controller.Success(g, body)
	}
}

// ResponseSuccessOnly
/**
* Return function response success only, without data
 */
func (controller *BaseController[T]) ResponseSuccessOnly(serviceFunc func(g *gin.Context) *common.Error) gin.HandlerFunc {
	return func(g *gin.Context) {
		err := serviceFunc(g)
		if err != nil {
			controller.Error(g, err, nil)
			return
		}
		controller.Success(g, nil)
	}
}
