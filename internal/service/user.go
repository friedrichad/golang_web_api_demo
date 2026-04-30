package service

import (
	"time"
	"strconv"
	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/dtos"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/repository"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
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
	var query dtos.UserFilter
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
	idStr := c.Param("id")
	if idStr == "" {
		return nil, common.RequestInvalid
	}

	// Convert string to int
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, common.NotFound
	}

	if user == nil {
		return nil, &common.Error{Code: "404", Message: "Người dùng không tồn tại"}
	}

	userResponse := modelToUserResponse(user)
	return &userResponse, nil
}

func (s *UserService) CreateUser(c *gin.Context) (*dtos.UserResponse, *common.Error) {
	var req dtos.UserCreate
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
	var req dtos.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		return common.RequestInvalid
	}

	// Convert UserID string to int
	userId, err := strconv.Atoi(req.UserID)
	if err != nil {
		return common.RequestInvalid
	}

	// Get existing user
	existingUser, err := s.userRepo.GetById(userId)
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
	if req.StatusInt != nil {
		existingUser.StatusInt = *req.StatusInt
	}
	if req.UpdatedBy != nil {
		existingUser.UpdatedBy = *req.UpdatedBy
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
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return common.RequestInvalid
		}
		ids[i] = id
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
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, common.RequestInvalid
	}

	authorities, err := s.userRepo.GetAuthorities(userId)
	if err != nil {
		return nil, common.SystemError
	}

	return authorities, nil
}


// Helper function to convert User model to UserResponse DTO
func modelToUserResponse(user *model.User) dtos.UserResponse {
	return dtos.UserResponse{
		UserID:      strconv.Itoa(user.UserID),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		StatusInt:   int(user.StatusInt),
		CreatedBy:   user.CreatedBy,
		CreatedAt:   user.CreatedAt,
		UpdatedBy:   user.UpdatedBy,
		UpdatedAt:   user.UpdatedAt,
	}
}
