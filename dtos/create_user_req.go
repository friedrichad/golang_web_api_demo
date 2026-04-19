package dtos

import (
	"net/mail"
	"strings"
	"fmt"
)

type UserRequest struct {

    Username    string `json:"username" binding:"required,alphanum,min=3,max=50"`
    DisplayName string `json:"display_name"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    RoleID      int    `json:"role_id" binding:"required,gt=0"`
}

func (r *UserRequest) Verify() (bool, error) {
	if r.RoleID < 1 {
		return false, fmt.Errorf("RoleId is required")
	}
	if strings.TrimSpace(r.Username) == "" || strings.TrimSpace(r.Email) == "" {
		return false, fmt.Errorf("Username or Email is empty")
	}
	if len(r.Password) < 6 {
		return false, fmt.Errorf("Password must be at least 6 characters")
	}

	addr, err := mail.ParseAddress(r.Email)
	if err != nil || addr.Address != r.Email {
		return false, fmt.Errorf("Email is invalid")
	}

	if strings.TrimSpace(r.DisplayName) == "" {
		r.DisplayName = r.Username
	}
	return true, nil
}
