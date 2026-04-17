package dtos

type UserResp struct {
	UserID      int32  `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	StatusInt   int32  `json:"status_int"`
}
