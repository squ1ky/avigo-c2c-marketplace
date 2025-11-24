package dto

type UserInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}

type ProfileResponse struct {
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone,omitempty"`
	About       string `json:"about,omitempty"`
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" validate:"omitempty,min=2,max=64"`
	Phone       *string `json:"phone" validate:"omitempty,max=32"`
	About       *string `json:"about" validate:"omitempty,max=2000"`
	Country     *string `json:"country" validate:"omitempty,max=64"`
	City        *string `json:"city" validate:"omitempty,max=64"`
	AvatarURL   *string `json:"avatar_url" validate:"omitempty,url,max=512"`
}
