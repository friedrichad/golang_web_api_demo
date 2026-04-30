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
	GetByUuid(id string) (*model.User, error)
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
	var users []model.User
	var total int64
	q := u.DB.Model(&model.User{})

	// Apply filters
	if query.Username != "" {
		q = q.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.DisplayName != "" {
		q = q.Where("display_name LIKE ?", "%"+query.DisplayName+"%")
	}
	if query.Email != "" {
		q = q.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.StatusInt != 0 {
		q = q.Where("status_int = ?", query.StatusInt)
	}

	// Get total count
	err := q.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (query.Page - 1) * query.Size
	err = q.Offset(offset).Limit(query.Size).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
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

func (u *UserRepository) GetByUuid(id string) (*model.User, error) {
	var user *model.User
	err := u.Instance.Where("user_uuid = ?", id).First(&user).Error
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
