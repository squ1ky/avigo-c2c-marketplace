package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	pgrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	"time"
)

var (
	ErrOrderNotCompleted   = errors.New("order is not completed")
	ErrReviewAlreadyExists = errors.New("review already exists for this order")
	ErrReviewForbidden     = errors.New("")
)

type ReviewService struct {
	reviewRepo *pgrepo.ReviewRepository
	orderRepo  *pgrepo.OrderRepository
}

func NewReviewService(
	reviewRepo *pgrepo.ReviewRepository,
	orderRepo *pgrepo.OrderRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, input dto.CreateReviewInput) (uuid.UUID, error) {
	order, err := s.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != domain.OrderStatusCompleted {
		return uuid.Nil, ErrOrderNotCompleted
	}
	if order.BuyerID != input.ReviewerID {
		return uuid.Nil, ErrReviewForbidden
	}

	exists, err := s.reviewRepo.ExistsByOrderID(ctx, input.OrderID)
	if err != nil {
		return uuid.Nil, nil
	}
	if exists {
		return uuid.Nil, ErrReviewAlreadyExists
	}

	reviewID := uuid.New()
	now := time.Now()

	review := &domain.Review{
		ID:             reviewID,
		OrderID:        input.OrderID,
		ReviewerID:     input.ReviewerID,
		ReviewedUserID: order.SellerID,
		Rating:         input.Rating,
		Text:           input.Text,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create review: %w", err)
	}

	return reviewID, nil
}

func (s *ReviewService) GetSellerReviews(ctx context.Context, sellerID uuid.UUID) ([]dto.ReviewResponse, error) {
	reviews, err := s.reviewRepo.GetByUserID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	result := make([]dto.ReviewResponse, len(reviews))
	for i, review := range reviews {
		result[i] = dto.ReviewResponse{
			ID:         review.ID,
			Rating:     review.Rating,
			Text:       review.Text,
			CreatedAt:  review.CreatedAt,
			ReviewerID: review.ReviewerID,
		}
	}

	return result, nil
}
