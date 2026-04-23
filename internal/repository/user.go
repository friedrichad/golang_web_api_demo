package repository

import(
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
)
type IUserRepository interface {
	IBaseRepository[model.User, int]
	GetByUsername(username string) (*model.User, error)
	GetAuthorities(userId int) ([]string, error)
	GetAllByCondition(query model.UserRequest) ([]model.User, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.User, error)
	Save(user *model.User) error
	Update(user *model.UserUpdate) error
}

type UserRepository struct {
	BaseRepository[model.User, int]
	DB *gorm.DB
}
var userRepository IUserRepository

func NewUserRepository() IUserRepository {
	if userRepository == nil {
		userRepository = &UserRepository{DB: database.Instance}
		userRepository.SetInstance(db.Instance)
		return userRepository
	}
	return userRepository
}