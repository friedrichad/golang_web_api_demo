package shared

type WelcomeMailData struct{
	UserID int `json:"user_id"`
	Email string `json:"email"`
}