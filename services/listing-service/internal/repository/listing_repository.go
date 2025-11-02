package repository

import (
	"context"

	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
)

type ListingRepository interface {
	Create(ctx context.Context, listing *domain.Listing) error
	GetByID(ctx context.Context, id string) (*domain.Listing, error)
	Update(ctx context.Context, listing *domain.Listing) error
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Listing, error)
	AddPhoto(ctx context.Context, photo *domain.Photo) error
	GetPhotos(ctx context.Context, listingID []string) ([]*domain.Photo, error)
	DeletePhoto(ctx context.Context, photoId string) error
}
