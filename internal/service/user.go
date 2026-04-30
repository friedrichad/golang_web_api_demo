package service

import (
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	GetAllUsers(c *gin.Context) ([]dtos.UserResponse, int, *common.Error)
	GetUserByUuid(c *gin.Context) (*dtos.UserResponse, *common.Error)
	CreateUser(c *gin.Context) (*dtos.UserResponse, *common.Error)
	UpdateUser(c *gin.Context) *common.Error
	DeleteUser(c *gin.Context) *common.Error
	GetUserAuthorities(c *gin.Context) ([]string, *common.Error)
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

func (s *UserService) GetAllUsers(c *gin.Context) ([]dtos.UserResponse, int, *common.Error) {
	var query dtos.UserRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

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
	if total == 0 {
		return nil, 0, common.NotFound
	}

	// Convert models to DTOs
	userResponses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = modelToUserResponse(&user)
	}

	return userResponses, total, nil
}

func (s *UserService) GetUserByUuid(c *gin.Context) (*dtos.UserResponse, *common.Error) {
	id := c.Param("id")
	if id == "" {
		return nil, common.RequestInvalid
	}

	user, err := s.userRepo.GetByUuid(id)
	if err != nil {
		return nil, common.NotFound
	}

	if user == nil {
		return nil, &common.Error{Code: "404", Message: "Người dùng không tồn tại"}
	}

	UserResponse := modelToUserResponse(user)
	return &UserResponse, nil
}

func (s *UserService) CreateUser(c *gin.Context) (*dtos.UserResponse, *common.Error) {
	var req dtos.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, common.RequestInvalid
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, &common.Error{Code: "400", Message: "Username đã tồn tại"}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, common.SystemError
	}

	// Create user model
	user := &model.User{
		UserUUID:     uuid.New().String(),
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		StatusInt:    1, // Active by default
		CreatedAt: time.Now(),
	}

	err = s.userRepo.Save(user)
	if err != nil {
		return nil, common.SystemError
	}

	UserResponse := modelToUserResponse(user)
	return &UserResponse, nil
}

func (s *UserService) UpdateUser(c *gin.Context) *common.Error {
	var user dtos.UserUpdate
	if err := c.ShouldBindJSON(&user); err != nil {
		return common.RequestInvalid
	}

	// Get existing user
	existingUser, err := s.userRepo.GetByUuid(user.UserUUID)
	if err != nil {
		return common.NotFound
	}

	if existingUser == nil {
		return &common.Error{Code: "404", Message: "Người dùng không tồn tại"}
	}

	// Update fields
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.DisplayName != "" {
		existingUser.DisplayName = user.DisplayName
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.StatusInt != nil {
		existingUser.StatusInt = *user.StatusInt
	}
	if user.UpdatedBy != nil {
		existingUser.UpdatedBy = *user.UpdatedBy
	}

	// Update password if provided
	if user.NewPassword != "" {
		// Hash new password
		hashedPassword, errHash := utils.HashPassword(user.NewPassword)
		if errHash != nil {
			return common.SystemError
		}
		existingUser.PasswordHash = hashedPassword
	}

	err = s.userRepo.Update(existingUser)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *UserService) DeleteUser(c *gin.Context) *common.Error {
	var ids []string
	if err := c.ShouldBindJSON(&ids); err != nil {
		return common.RequestInvalid
	}

	err := s.userRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *UserService) GetUserAuthorities(c *gin.Context) ([]string, *common.Error) {
	id := c.Param("id")
	if id == "" {
		return nil, common.RequestInvalid
	}

	authorities, err := s.userRepo.GetAuthorities(id)
	if err != nil {
		return nil, common.SystemError
	}

	return authorities, nil
}

// Helper function to convert User model to UserResponse DTO
func modelToUserResponse(user *model.User) dtos.UserResponse {
	return dtos.UserResponse{
		UserUUID:    user.UserUUID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		StatusInt:   user.StatusInt,
		CreatedBy:   user.CreatedBy,
		CreatedAt:   user.CreatedAt,
		UpdatedBy:   user.UpdatedBy,
		UpdatedAt:   user.UpdatedAt,
	}
}
