package model

import (
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"time"
)

type DateRequest struct {
	DateFrom *time.Time `form:"dateFrom" time_format:"2006-01-02 15:04:05"`
	DateTo   *time.Time `form:"dateTo" time_format:"2006-01-02 15:04:05"`
}

func (d DateRequest) GetDateFrom() *time.Time {
	return utils.TrunDate(d.DateFrom, true)
}

func (d DateRequest) GetDateTo() *time.Time {
	return utils.TrunDate(d.DateTo, false)
}

type PageSize struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
