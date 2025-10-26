package dto

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,username"`
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,password"`
	DisplayName string `json:"display_name" validate:"required,displayname"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,max=32"`
	Country     string `json:"country,omitempty" validate:"omitempty,max=64"`
	City        string `json:"city,omitempty" validate:"omitempty,max=64"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required,min=3"`
	Password   string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

type ConfirmEmailRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Code   string `json:"code" validate:"required,len=6"`
}

type UserInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}
