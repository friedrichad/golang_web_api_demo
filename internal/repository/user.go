package repository

import (
	"github.com/friedrichad/golang_web_api_demo/internal/configs/db"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	IBaseRepository[model.User, int]
	GetByUsername(username string) (*model.User, error)
	GetAuthorities(userId int) ([]string, error)
	GetAllByCondition(query model.UserRequest) ([]model.User, int, error)
	Delete(ids []int) error
	GetByUuid(id int) (*model.User, error)
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

func (u *UserRepository) GetAllByCondition(query model.UserRequest) ([]model.User, int, error) {
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

func (u *UserRepository) GetByUuid(id int) (*model.User, error) {
	var user *model.User
	err := u.Instance.Where("user_uuid = ?", id).First(&user).Error
	return user, err
}

func (u *UserRepository) Save(user *model.User) error {
	return u.Instance.Create(user).Error
}

func (u *UserRepository) Update(user *model.UserUpdate) error {
	var existingUser *model.User
	err := u.Instance.Where("user_uuid = ?", user.UserUUID).First(&existingUser).Error
	if err != nil {
		return err
	}
	existingUser.Username = user.Username
	existingUser.DisplayName = user.DisplayName
	existingUser.Email = user.Email
	if user.NewPassword != "" {
		existingUser.PasswordHash = user.NewPassword // In production, hash the password
	}
	err = u.BaseRepository.Update(existingUser)
	return err
}

func (u *UserRepository) Delete(ids []int) error {
	err := u.BaseRepository.Delete(ids)
	return err
}
