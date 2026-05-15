package service

import (
	"strconv"
	"time"

	"github.com/friedrichad/golang_web_api_demo/internal/common"	
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetAllUsers(c *gin.Context) ([]model.UserResponse, int, *common.Error)
	GetUserByUuid(c *gin.Context) (*model.UserResponse, *common.Error)
	CreateUser(c *gin.Context) (*model.UserResponse, *common.Error)
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

func (s *UserService) GetAllUsers(c *gin.Context) ([]model.UserResponse, int, *common.Error) {
	var query model.UserFilter
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, 0, common.RequestInvalid
	}

	users, total, err := s.userRepo.GetAllByCondition(query)
	if err != nil {
		return nil, 0, common.SystemError
	}
	if total == 0 {
		return nil, 0, common.NotFound
	}

	// Convert models to DTOs
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = modelToUserResponse(&user)
	}

	return userResponses, total, nil
}

func (s *UserService) GetUserByUuid(c *gin.Context) (*model.UserResponse, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	// Convert string to int
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	user, err := s.userRepo.GetById(int(userId))
	if err != nil {
		return nil, common.NotFound
	}

	if user == nil {
		return nil, &common.Error{Code: "404", Message: "Người dùng không tồn tại"}
	}

	userResponse := modelToUserResponse(user)
	return &userResponse, nil
}

func (s *UserService) CreateUser(c *gin.Context) (*model.UserResponse, *common.Error) {
	var req model.UserCreate
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
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		StatusInt:    1, // Active by default
		CreatedAt:    time.Now(),
	}

	err = s.userRepo.Save(user)
	if err != nil {
		return nil, common.SystemError
	}

	userResponse := modelToUserResponse(user)
	return &userResponse, nil
}

func (s *UserService) UpdateUser(c *gin.Context) *common.Error {
	var req model.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	// Get existing user
	existingUser, err := s.userRepo.GetById(req.UserID)
	if err != nil {
		return common.NotFound
	}

	if existingUser == nil {
		return &common.Error{Code: "404", Message: "Người dùng không tồn tại"}
	}

	// Update fields
	if req.Username != "" {
		existingUser.Username = req.Username
	}
	if req.DisplayName != "" {
		existingUser.DisplayName = req.DisplayName
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.StatusInt != 0 {
		existingUser.StatusInt = req.StatusInt
	}
	if req.UpdatedBy != 0 {
		existingUser.UpdatedBy = req.UpdatedBy
	}

	// Update password if provided
	if req.NewPassword != "" {
		// Hash new password
		hashedPassword, errHash := utils.HashPassword(req.NewPassword)
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
	var idStrs []string
	if err := c.ShouldBindJSON(&idStrs); err != nil {
		return common.RequestInvalid
	}

	// Convert string IDs to int
	ids := make([]int, len(idStrs))
	for i, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = int(id)
	}

	err := s.userRepo.Delete(ids)
	if err != nil {
		return common.SystemError
	}

	return nil
}

func (s *UserService) GetUserAuthorities(c *gin.Context) ([]string, *common.Error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	// Convert string to int
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, common.RequestInvalid
	}

	authorities, err := s.userRepo.GetAuthorities(int(userId))
	if err != nil {
		return nil, common.SystemError
	}

	return authorities, nil
}

// Helper function to convert User model to UserResponse DTO
func modelToUserResponse(user *model.User) model.UserResponse {
	return model.UserResponse{
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		StatusInt:   int(user.StatusInt),
		CreatedBy:   int(user.CreatedBy),
		CreatedAt:   user.CreatedAt,
		UpdatedBy:   int(user.UpdatedBy),
		UpdatedAt:   user.UpdatedAt,
	}
}
