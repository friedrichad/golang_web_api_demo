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
