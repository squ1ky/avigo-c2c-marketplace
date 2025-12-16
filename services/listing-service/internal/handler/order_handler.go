package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"net/http"
)

type CreateOrderRequest struct {
	ListingID uuid.UUID `json:"listing_id" binding:"required"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/account/:user_id/orders")
	{
		orders.GET("/purchases", h.getUserPurchases)
		orders.GET("/sales", h.getUserSales)
	}

	orderActions := api.Group("/orders")
	{
		orderActions.POST("", h.createOrder)
		orderActions.POST("/:order_id/confirm", h.confirmOrder)
		orderActions.POST("/:order_id/complete", h.completeOrder)
		orderActions.POST("/:order_id/cancel", h.cancelOrder)
	}
}

func (h *Handler) getUserPurchases(c *gin.Context) {
	targetUserID, err := h.getTargetUserID(c)
	if err != nil {
		badRequest(c, "INVALID_USER_ID", err)
		return
	}

	requesterID, ok := getUserIDFromHeader(c)
	if !ok || requesterID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	resp, err := h.orderSvc.GetUserOrders(c.Request.Context(), targetUserID, false)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getUserSales(c *gin.Context) {
	targetUserID, err := h.getTargetUserID(c)
	if err != nil {
		badRequest(c, "INVALID_USER_ID", err)
		return
	}

	requesterID, ok := getUserIDFromHeader(c)
	if !ok || requesterID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	resp, err := h.orderSvc.GetUserOrders(c.Request.Context(), targetUserID, true)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) createOrder(c *gin.Context) {
	buyerID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidJSON(c, err)
		return
	}

	orderID, err := h.orderSvc.CreateOrder(c.Request.Context(), buyerID, req.ListingID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order_id": orderID,
		"status":   domain.OrderStatusCreated,
	})
}

func (h *Handler) confirmOrder(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		badRequest(c, "INVALID_ORDER_ID", err)
		return
	}

	if err := h.orderSvc.ConfirmOrder(c.Request.Context(), orderID, userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) completeOrder(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		badRequest(c, "INVALID_ORDER_ID", err)
		return
	}

	if err := h.orderSvc.CompleteOrder(c.Request.Context(), orderID, userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) cancelOrder(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		badRequest(c, "INVALID_ORDER_ID", err)
		return
	}

	var req CancelOrderRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.orderSvc.CancelOrder(c.Request.Context(), orderID, userID, req.Reason); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
