package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/grpc/client/user"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/mapper"
	pgrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	userpb "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/pb/user"
	"log/slog"
	"time"
)

var (
	ErrListingNotFound  = errors.New("listing not found")
	ErrListingNotActive = errors.New("listing not active")
	ErrListingSold      = errors.New("listing is already sold")
	ErrSelfPurchase     = errors.New("cannot buy your own listing")

	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid order status")

	ErrForbidden = errors.New("forbidden action")
)

// OrderService contains business logic for marketplace orders and deal lifecycle.
type OrderService struct {
	orderRepo   *pgrepo.OrderRepository
	listingRepo *pgrepo.ListingRepository
	txManager   pgrepo.TransactionManager
	userClient  *user.Client
}

func NewOrderService(
	orderRepo *pgrepo.OrderRepository,
	listingRepo *pgrepo.ListingRepository,
	txManager pgrepo.TransactionManager,
	userClient *user.Client,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		listingRepo: listingRepo,
		txManager:   txManager,
		userClient:  userClient,
	}
}

// CreateOrder creates a new order for the given listing and buyer,
// marks the listing as sold inside a DB transaction, and returns the order ID.
func (s *OrderService) CreateOrder(ctx context.Context, buyerID uuid.UUID, listingID uuid.UUID) (uuid.UUID, error) {
	listing, err := s.listingRepo.GetByID(ctx, listingID)
	if err != nil {
		return uuid.Nil, ErrListingNotFound
	}

	if listing.Status != domain.ListingStatusActive {
		return uuid.Nil, ErrListingNotActive
	}
	if listing.IsSold {
		return uuid.Nil, ErrListingSold
	}
	if listing.UserID == buyerID {
		return uuid.Nil, ErrSelfPurchase
	}

	orderID := uuid.New()
	now := time.Now()

	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.listingRepo.MarkAsSold(txCtx, listingID); err != nil {
			return fmt.Errorf("failed to mark listing sold: %w", err)
		}

		order := &domain.Order{
			ID:        orderID,
			ListingID: listingID,
			BuyerID:   buyerID,
			SellerID:  listing.UserID,
			Status:    domain.OrderStatusCreated,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := s.orderRepo.Create(txCtx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

// ConfirmOrder marks the order as confirmed by the seller,
// allowed only for the seller and only from the "created" status.
func (s *OrderService) ConfirmOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if order.SellerID != userID {
		return ErrForbidden
	}

	if order.Status != domain.OrderStatusCreated {
		return ErrInvalidStatus
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusConfirmed, "")
}

// CompleteOrder marks the order as completed by the buyer,
// allowed only for the buyer and only from the "confirmed" status.
func (s *OrderService) CompleteOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if order.BuyerID != userID {
		return ErrForbidden
	}

	if order.Status != domain.OrderStatusConfirmed {
		return ErrInvalidStatus
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusCompleted, "")
}

// CancelOrder cancels the order with an optional reason,
// allowed for buyer or seller, and returns the listing back to "active" status in a transaction.
func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, reason string) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return ErrForbidden
	}

	if order.Status == domain.OrderStatusCompleted || order.Status == domain.OrderStatusCancelled {
		return ErrInvalidStatus
	}

	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.orderRepo.UpdateStatus(txCtx, orderID, domain.OrderStatusCancelled, reason); err != nil {
			return err
		}

		if err := s.listingRepo.MarkAsActive(txCtx, order.ListingID); err != nil {
			return fmt.Errorf("failed to return listing to active: %w", err)
		}

		return nil
	})
}

// GetUserOrders returns a list of orders for the given user:
// purchases if isSeller == false, or sales if isSeller == true
// enriched with counterparty user data fetched via gRPC.
func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID, isSeller bool) ([]dto.OrderResponse, error) {
	var orders []domain.Order
	var err error

	if isSeller {
		orders, err = s.orderRepo.GetSalesByUserID(ctx, userID)
	} else {
		orders, err = s.orderRepo.GetPurchasesByUserID(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	if len(orders) == 0 {
		return []dto.OrderResponse{}, nil
	}

	userIDs := make([]string, 0, len(orders))
	for _, order := range orders {
		otherID := order.SellerID
		if isSeller {
			otherID = order.BuyerID
		}
		userIDs = append(userIDs, otherID.String())
	}

	usersMap, err := s.userClient.GetUsersByID(ctx, userIDs)
	if err != nil {
		slog.Warn("warning: user service unavailable: %v\n", err)
		usersMap = make(map[string]*userpb.User)
	}

	responses := mapper.OrdersToResponses(orders, usersMap, isSeller)

	return responses, nil
}
