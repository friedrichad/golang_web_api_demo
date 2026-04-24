package service

import (
	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	GetAllUsers(c *gin.Context, query model.UserRequest) ([]model.User, int, *common.Error)
	GetUserById(c *gin.Context, id int) (*model.User, *common.Error)
	CreateUser(c *gin.Context, user *model.User) (*model.User, *common.Error)
	UpdateUser(c *gin.Context, user *model.UserUpdate) *common.Error
	DeleteUser(c *gin.Context, ids []int) *common.Error
	GetUserAuthorities(c *gin.Context, userId int) ([]string, *common.Error)
}

type UserService struct {
	userRepo repository.IUserRepository
}

var userService IUserService

func NewUserService() IUserService {
	if userService == nil {
		userService = &UserService{
			userRepo: repository.NewUserRepository(),
		}
	}
	return userService
}

func (s *UserService) GetAllUsers(c *gin.Context, query model.UserRequest) ([]model.User, int, *common.Error) {
	// Validate pagination
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}

	users, total, err := s.userRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}

	return users, total, nil
}

func (s *UserService) GetUserById(c *gin.Context, id int) (*model.User, *common.Error) {
	user, err := s.userRepo.GetById(id)
	if err != nil {
		return nil, common.SystemError
	}

	if user == nil {
		return nil, &common.Error{Code: "404", Message: "Người dùng không tồn tại"}
	}

	return user, nil
}

func (s *UserService) CreateUser(c *gin.Context, user *model.User) (*model.User, *common.Error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(user.Username)
	if existingUser != nil {
		return nil, &common.Error{Code: "400", Message: "Username đã tồn tại"}
	}

	// Generate UUID
	user.UserUUID = uuid.New().String()

	err = s.userRepo.Save(user)
	if err != nil {
		return nil, common.SystemError
	}

	return user, nil
}

func (s *UserService) UpdateUser(c *gin.Context, user *model.UserUpdate) *common.Error {
	err := s.userRepo.Update(user)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *UserService) DeleteUser(c *gin.Context, ids []int) *common.Error {
	err := s.userRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *UserService) GetUserAuthorities(c *gin.Context, userId int) ([]string, *common.Error) {
	authorities, err := s.userRepo.GetAuthorities(userId)
	if err != nil {
		return nil, common.SystemError
	}

	return authorities, nil
}
