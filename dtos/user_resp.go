package dtos

type UserResponse struct {
	UserID      int32  `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	StatusInt   int32  `json:"status_int"`
	RoleID      int32  `json:"role_id"`
}
