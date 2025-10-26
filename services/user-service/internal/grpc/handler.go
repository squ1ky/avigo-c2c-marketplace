package grpc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/user-service/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	userRepo repository.UserRepository
}

func NewUserServiceServer(userRepo repository.UserRepository) *UserServiceServer {
	return &UserServiceServer{userRepo: userRepo}
}

func (s *UserServiceServer) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &userpb.GetUserByIDResponse{
		User: s.convertToProtoUser(user),
	}, nil
}

func (s *UserServiceServer) GetUsersByIDs(ctx context.Context, req *userpb.GetUsersByIDsRequest) (*userpb.GetUsersByIDsResponse, error) {
	users := make([]*userpb.User, 0, len(req.UserIds))

	for _, idStr := range req.UserIds {
		userID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			continue
		}

		users = append(users, s.convertToProtoUser(user))
	}

	return &userpb.GetUsersByIDsResponse{
		Users: users,
	}, nil
}

func (s *UserServiceServer) ValidateUser(ctx context.Context, req *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &userpb.ValidateUserResponse{
			Exists:   false,
			IsActive: false,
		}, nil
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return &userpb.ValidateUserResponse{
			Exists:   false,
			IsActive: false,
		}, nil
	}

	return &userpb.ValidateUserResponse{
		Exists:   true,
		IsActive: user.Status == domain.StatusActive,
	}, nil
}

func (s *UserServiceServer) convertToProtoUser(user *domain.User) *userpb.User {
	protoUser := &userpb.User{
		Id:       user.ID.String(),
		Username: user.Username,
		Role:     string(user.Role),
		Status:   string(user.Status),
	}

	if user.Security != nil {
		protoUser.Email = user.Security.Email
	}

	if user.Profile != nil {
		protoUser.DisplayName = user.Profile.DisplayName
		protoUser.Profile = &userpb.UserProfile{
			Phone:     user.Profile.Phone,
			About:     user.Profile.About,
			Country:   user.Profile.Country,
			City:      user.Profile.City,
			AvatarUrl: user.Profile.AvatarURL,
		}
	}

	return protoUser
}
