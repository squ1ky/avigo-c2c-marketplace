package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID          uuid.UUID            `json:"id"`
	Username    string               `json:"username"`
	Email       string               `json:"email"`
	DisplayName string               `json:"display_name"`
	Role        string               `json:"role"`
	Status      string               `json:"status"`
	Profile     *UserProfileResponse `json:"profile,omitempty"`
}

type UserProfileResponse struct {
	Phone     string `json:"phone,omitempty"`
	About     string `json:"about,omitempty"`
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
