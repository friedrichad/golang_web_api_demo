package repository
// import (
// 	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
// 	"github.com/friedrichad/golang_web_api_demo/internal/model"
// 	"gorm.io/gorm"
// )

// type IBinRespository interface {
// 	IBaseRepository[model.User, int]
// 	GetByUsername(username string) (*model.User, error)
// 	GetAuthorities(userId int) ([]string, error)
// 	GetAllByCondition(query model.UserRequest) ([]model.User, int, error)
// 	Delete(ids []int) error
// 	GetById(id int) (*model.User, error)
// 	Save(user *model.User) error
// 	Update(user *model.UserUpdate) error
// }

// type BinRespository struct {
// 	BaseRepository[model.User, int]
// 	DB *gorm.DB
// }

// var binRepository IBinRespository