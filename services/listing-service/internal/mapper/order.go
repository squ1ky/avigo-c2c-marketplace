package mapper

import (
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/pb/user"
)

func OrdersToResponses(
	orders []domain.Order,
	users map[string]*userpb.User,
	isSeller bool,
) []dto.OrderResponse {
	res := make([]dto.OrderResponse, len(orders))

	for i, order := range orders {
		var counterpartID uuid.UUID
		if isSeller {
			counterpartID = order.BuyerID
		} else {
			counterpartID = order.SellerID
		}

		u := users[counterpartID.String()]

		res[i] = dto.OrderResponse{
			ID:        order.ID,
			ListingID: order.ListingID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,

			CounterpartyID:   counterpartID,
			CounterpartyName: safeUsername(u),
			CounterpartyImg:  safeUserAvatar(u),
		}
	}

	return res
}

func safeUsername(u *userpb.User) string {
	if u == nil {
		return "Unknown User"
	}
	if u.DisplayName != "" {
		return u.DisplayName
	}
	if u.Username != "" {
		return u.Username
	}
	return "User"
}

func safeUserAvatar(u *userpb.User) string {
	if u != nil && u.Profile != nil {
		return u.Profile.AvatarUrl
	}
	return ""
}
