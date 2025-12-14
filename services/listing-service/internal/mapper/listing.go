package mapper

import (
	"fmt"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/pb/user"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
)

func ToDomainListing(input dto.CreateListingInput, id uuid.UUID, now time.Time) *domain.Listing {
	return &domain.Listing{
		ID:          id,
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Currency:    input.Currency,
		Status:      domain.ListingStatusActive,
		IsSold:      false,
		ViewsCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func ToDomainCharacteristics(input dto.CreateListingInput, id uuid.UUID, now time.Time) *domain.ListingCharacteristics {
	chars := input.Characteristics
	if chars == nil {
		chars = make(map[string]interface{})
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	return &domain.ListingCharacteristics{
		ListingID:       id,
		Characteristics: chars,
		Tags:            tags,
		UpdatedAt:       now,
	}
}

func ToListingResponse(
	listing *domain.Listing,
	chars *domain.ListingCharacteristics,
	media []domain.ListingMedia,
	baseURL string,
) *dto.ListingResponse {
	mediaResp := make([]dto.ListingMediaResponse, 0, len(media))

	for _, m := range media {
		fullURL := m.FileURL

		if baseURL != "" && !strings.HasPrefix(fullURL, "http") {
			fullURL = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(m.FileURL, "/"))
		}

		mediaResp = append(mediaResp, dto.ListingMediaResponse{
			ID:       m.ID,
			FileURL:  fullURL,
			FileType: m.FileType,
			Order:    m.Order,
		})
	}

	resp := &dto.ListingResponse{
		ID:              listing.ID,
		UserID:          listing.UserID,
		CategoryID:      listing.CategoryID,
		Title:           listing.Title,
		Description:     listing.Description,
		Price:           listing.Price,
		Currency:        listing.Currency,
		Status:          listing.Status,
		ViewsCount:      listing.ViewsCount,
		IsSold:          listing.IsSold,
		CreatedAt:       listing.CreatedAt,
		UpdatedAt:       listing.UpdatedAt,
		Media:           mediaResp,
		Characteristics: make(map[string]interface{}),
		Tags:            make([]string, 0),
	}

	if chars != nil {
		if chars.Characteristics != nil {
			resp.Characteristics = chars.Characteristics
		}
		if chars.Tags != nil {
			resp.Tags = chars.Tags
		}
	}

	return resp
}

func ToListingWithUserResponse(
	listing *domain.Listing,
	chars *domain.ListingCharacteristics,
	media []domain.ListingMedia,
	user *userpb.User,
	baseURL string,
) *dto.ListingWithUserResponse {

	listingResp := ToListingResponse(listing, chars, media, baseURL)

	var userResp dto.UserResponse

	if user != nil {
		userID, _ := uuid.Parse(user.Id)

		userResp = dto.UserResponse{
			ID:          userID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Status:      user.Status,
			Role:        user.Role,
		}

		if user.Profile != nil {
			userResp.Profile = &dto.UserProfileResponse{
				Phone:     user.Profile.Phone,
				About:     user.Profile.About,
				Country:   user.Profile.Country,
				City:      user.Profile.City,
				AvatarURL: user.Profile.AvatarUrl,
			}
		}
	} else {
		userResp.ID = listing.UserID
	}

	return &dto.ListingWithUserResponse{
		Listing: *listingResp,
		User:    userResp,
	}
}
