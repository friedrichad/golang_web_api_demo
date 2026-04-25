package repository

// import (
// 	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
// 	"github.com/friedrichad/golang_web_api_demo/internal/model"
// 	"gorm.io/gorm"
// )

// type IRequestDetail interface {
// 	IBaseRepository[model.User, int]
// 	GetByUsername(username string) (*model.User, error)
// 	GetAuthorities(userId int) ([]string, error)
// 	GetAllByCondition(query model.UserRequest) ([]model.User, int, error)
// 	Delete(ids []int) error
// 	GetByUuid(id int) (*model.User, error)
// 	Save(user *model.User) error
// 	Update(user *model.UserUpdate) error
// }