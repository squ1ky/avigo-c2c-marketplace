package mapper

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/pb/user"
	"strings"
)

func OrdersToResponses(
	orders []domain.Order,
	usersMap map[string]*userpb.User,
	isSeller bool,
	listingsMap map[uuid.UUID]*domain.Listing,
	listingMediaMap map[uuid.UUID][]domain.ListingMedia,
	baseURL string,
) []dto.OrderResponse {
	res := make([]dto.OrderResponse, len(orders))

	for i, order := range orders {
		var counterpartID uuid.UUID
		if isSeller {
			counterpartID = order.BuyerID
		} else {
			counterpartID = order.SellerID
		}

		userResp := dto.UserResponse{
			ID:          counterpartID,
			Username:    "unknown",
			DisplayName: "Неизвестный",
		}

		if u, ok := usersMap[counterpartID.String()]; ok {
			if uid, err := uuid.Parse(u.Id); err == nil {
				userResp.ID = uid
			}

			userResp.Username = u.Username
			userResp.Email = u.Email
			userResp.DisplayName = u.DisplayName
			userResp.Role = u.Role
			userResp.Status = u.Status

			if u.Profile != nil {
				userResp.Profile = &dto.UserProfileResponse{
					Phone:     u.Profile.Phone,
					About:     u.Profile.About,
					Country:   u.Profile.Country,
					City:      u.Profile.City,
					AvatarURL: u.Profile.AvatarUrl,
				}
			}
		}

		listingResp := dto.ListingShortResponse{
			ID:       order.ListingID,
			Title:    "Товар недоступен",
			Price:    0,
			Currency: "RUB", // Default
			Media:    []dto.ListingMediaResponse{},
		}

		if lst, ok := listingsMap[order.ListingID]; ok {
			listingResp.Title = lst.Title
			listingResp.Price = lst.Price
			listingResp.Currency = lst.Currency

			if mediaList, ok := listingMediaMap[order.ListingID]; ok && len(mediaList) > 0 {
				listingResp.Media = make([]dto.ListingMediaResponse, len(mediaList))

				for j, m := range mediaList {
					fullURL := m.FileURL
					if baseURL != "" && !strings.HasPrefix(fullURL, "http") {
						fullURL = fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(m.FileURL, "/"))
					}

					listingResp.Media[j] = dto.ListingMediaResponse{
						ID:       m.ID,
						FileURL:  fullURL,
						FileType: m.FileType,
						Order:    m.Order,
					}
				}
			}
		}

		res[i] = dto.OrderResponse{
			ID:           order.ID,
			Status:       order.Status,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
			Listing:      listingResp,
			Counterparty: userResp,
		}
	}

	return res
}
