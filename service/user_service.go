package service

import (
	"fmt"

	dtos "github.com/friedrichad/golang_web_api_demo/dtos"
	"github.com/friedrichad/golang_web_api_demo/models"
	"gorm.io/gorm"
)

type IUserService interface {
	GetUserByID(userID int32) (*dtos.UserResp, error)
	GetUsers() ([]*dtos.UserResp, error)
	PostUser(user *dtos.UserResp) (*dtos.UserResp, error)
}

type UserService struct {
	DB *gorm.DB
}

// GetUserByID retrieves a single user by their ID
func (s *UserService) GetUserByID(userID int32) (*dtos.UserResp, error) {
	var user models.User

	result := s.DB.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, result.Error
	}

	// Convert models.User to dtos.UserResp (hide password_hash)
	return &dtos.UserResp{
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		StatusInt:   user.StatusInt,
	}, nil
}

// GetUsers retrieves all users from the database
func (s *UserService) GetUsers() ([]*dtos.UserResp, error) {
	var users []models.User

	result := s.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert models.User to dtos.UserResp (hide password_hash)
	var userResps []*dtos.UserResp
	for _, user := range users {
		userResps = append(userResps, &dtos.UserResp{
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
func (s *UserService) PostUser(user *dtos.UserResp) (*dtos.UserResp, error) {
	// Validation
	if user.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Convert dtos.UserResp to models.User
	// Note: Password hash should be provided and properly hashed before calling this method
	dbUser := models.User{
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		PasswordHash: "", // This should be set by the handler with properly hashed password
		StatusInt:    user.StatusInt,
	}

	result := s.DB.Create(&dbUser)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert back to dtos.UserResp and return (without password hash)
	return &dtos.UserResp{
		UserID:      dbUser.UserID,
		Username:    dbUser.Username,
		DisplayName: dbUser.DisplayName,
		Email:       dbUser.Email,
		StatusInt:   dbUser.StatusInt,
	}, nil
}
