package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	IBaseRepository[model.User, int]
	GetByUsername(username string) (*model.User, error)
	GetAuthorities(userId int) ([]string, error)
	GetAllByCondition(query dtos.UserFilter) ([]model.User, int, error)
	Delete(ids []int) error
	GetById(id int) (*model.User, error)
	Save(user *model.User) error
	Update(user *model.User) error
}

type UserRepository struct {
	BaseRepository[model.User, int]
	DB *gorm.DB
}

var userRepository IUserRepository

func NewUserRepository() IUserRepository {
	if userRepository == nil {
		userRepository = &UserRepository{DB: db.Instance}
		userRepository.SetInstance(db.Instance)
	}
	return userRepository
}

func (u *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user *model.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (u *UserRepository) GetAllByCondition(query dtos.UserFilter) ([]model.User, int, error) {
	var username interface{}
	if query.Username != nil {
		username = *query.Username
	}

	var statusInt interface{}
	if query.StatusInt != nil {
		statusInt = *query.StatusInt
	}

	dateFrom := query.GetDateFrom()
	dateTo := query.GetDateTo()

	return u.GetPage(
		"Select u.* from user u"+
			" where (? is null or u.username = ?)"+
			" and (? is null or u.status_int = ?)"+
			" and (? is null or u.created_at >= ?)"+
			" and (? is null or u.created_at < ?)",
		query.Page,
		query.Size,
		username, username,
		statusInt, statusInt,
		dateFrom, dateFrom,
		dateTo, dateTo,
	)
}
func (u *UserRepository) GetAuthorities(userId int) ([]string, error) {
	var authorities []string
	err := u.Instance.Model(&model.Role{}).
		Select("role.role_name").
		Joins("INNER JOIN user_role ON role.role_id = user_role.role_id").
		Where("user_role.user_id = ?", userId).
		Pluck("role_name", &authorities).Error
	return authorities, err
}

func (u *UserRepository) GetById(id int) (*model.User, error) {
	var user *model.User
	err := u.Instance.Where("user_id = ?", id).First(&user).Error
	return user, err
}

func (u *UserRepository) Save(user *model.User) error {
	return u.BaseRepository.Create(user)
}

func (u *UserRepository) Update(user *model.User) error {
	return u.BaseRepository.Update(user)
}

func (u *UserRepository) Delete(ids []int) error {
	err := u.DB.Where("user_id IN ?", ids).Delete(&model.User{}).Error
	return err
}

func (u *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository[model.User, int]{Instance: tx},
		DB:             tx,
	}
}
