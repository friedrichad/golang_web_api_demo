package service

import (
	"fmt"

	dtos "github.com/friedrichad/golang_web_api_demo/dtos"
	"github.com/friedrichad/golang_web_api_demo/models"
	utils "github.com/friedrichad/golang_web_api_demo/utils"
	"gorm.io/gorm"
)

type IUserService interface {
	GetUserByID(userID int32) (*dtos.UserResponse, error)
	GetUsers() ([]*dtos.UserResponse, error)
	CreateUser(user *dtos.UserRequest) (*dtos.UserResponse, error)
}

type UserService struct {
	DB *gorm.DB
}

// GetUserByID retrieves a single user by their ID
func (s *UserService) GetUserByID(userID int32) (*dtos.UserResponse, error) {
	var user models.User

	result := s.DB.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, result.Error
	}

	// Convert models.User to dtos.UserResp (hide password_hash)
	return &dtos.UserResponse{
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		StatusInt:   user.StatusInt,
	}, nil
}

// GetUsers retrieves all users from the database
func (s *UserService) GetUsers() ([]*dtos.UserResponse, error) {
	var users []models.User

	result := s.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert models.User to dtos.UserResp (hide password_hash)
	var userResps []*dtos.UserResponse
	for _, user := range users {
		userResps = append(userResps, &dtos.UserResponse{
			UserID:      user.UserID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			StatusInt:   user.StatusInt,
		})
	}

	return userResps, nil
}

// PostUser creates a new user in the database
func (s *UserService) CreateUser(user *dtos.UserRequest) (*dtos.UserResponse, error) {
	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password: %v", err)
	}
	statusInt := int32(1)
	dbUser := models.User{
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		PasswordHash: passwordHash,
		StatusInt:    statusInt,
	}

	result := s.DB.Create(&dbUser)
	if result.Error != nil {
		return nil, result.Error
	}
	userRole := models.UR{
		UserID: dbUser.UserID,
		RoleID: int32(user.RoleID),
	}
	resultRole := s.DB.Create(&userRole)
	if resultRole.Error != nil {
		return nil, resultRole.Error
	}

	// Convert back to dtos.UserResponse and return (without password hash)
	return &dtos.UserResponse{
		UserID:      dbUser.UserID,
		Username:    dbUser.Username,
		DisplayName: dbUser.DisplayName,
		Email:       dbUser.Email,
		StatusInt:   dbUser.StatusInt,
		RoleID: userRole.RoleID,
	}, nil
}
