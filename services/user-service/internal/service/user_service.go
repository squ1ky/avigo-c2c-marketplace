package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/validation"
)

type UserService struct {
	userRepo  repository.UserRepository
	validator *validation.Validator
}

func NewUserService(userRepo repository.UserRepository, validator *validation.Validator) *UserService {
	return &UserService{
		userRepo:  userRepo,
		validator: validator,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Profile == nil {
		return nil, domain.ErrUserNotFound
	}

	p := user.Profile
	return &dto.ProfileResponse{
		DisplayName: p.DisplayName,
		Phone:       p.Phone,
		About:       p.About,
		Country:     p.Country,
		City:        p.City,
		AvatarURL:   p.AvatarURL,
		LastLoginAt: user.Security.LastLoginAt,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	updates := buildProfileUpdates(req)
	if len(updates) == 0 {
		return s.GetProfile(ctx, userID) // could be error here
	}

	if err := s.userRepo.UpdateProfileFields(ctx, userID, updates); err != nil {
		return nil, err
	}

	return s.GetProfile(ctx, userID)
}

func buildProfileUpdates(req dto.UpdateProfileRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.About != nil {
		updates["about"] = *req.About
	}
	if req.Country != nil {
		updates["country"] = *req.Country
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	return updates
}
