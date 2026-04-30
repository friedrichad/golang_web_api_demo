package service

import(
	"github.com/gin-gonic/gin"
	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
)

type IBinService interface {
	GetAllBins(c *gin.Context) ([]dtos.BinResponse, int, *common.Error)
	GetBinByBinId(c *gin.Context) (*dtos.BinResponse, *common.Error)
	CreateBin(c *gin.Context) (*dtos.BinResponse, *common.Error)
	UpdateBin(c *gin.Context) *common.Error
	DeleteBin(c *gin.Context) *common.Error
}


type BinService struct {
	binRepo repository.IBinRepository
}

var binService IBinService

func NewBinService() IBinService {
	if binService == nil {
		binService = &BinService{
			binRepo: repository.NewBinRepository(),
		}
	}
	return binService
}

// func (s *BinService) GetAllBins(c *gin.Context) ([]dtos.BinResponse, int, *common.Error) {
// 	var query dtos.Bin
// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		return nil, 0, common.RequestInvalid
// 	}
	