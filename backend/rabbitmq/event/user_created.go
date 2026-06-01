package event

type UserCreatedEvent struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func NewUserCreatedEvent(userID int, email string) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID: userID,
		Email:  email,
	}
}
