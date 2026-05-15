package model

type TokenInfo struct {
	Id          string   `json:"id"`
	Username    string   `json:"user_name"`
	Exp         int64    `json:"exp"`
	Authorities []string `json:"authorities"`
	ClientId    string   `json:"client_id"`
	Active      bool     `json:"active"`
}

type TokenResponse struct {
	TokenInfo
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	RefreshExp  int64  `json:"refresh_exp"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ClientID string `json:"client_id" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	Email       string `json:"email" binding:"required,email"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
	PositionID  int    `json:"position_id"`
}
