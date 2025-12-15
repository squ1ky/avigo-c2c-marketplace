package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/account/:user_id/orders")
	{
		orders.GET("/purchases", h.getUserPurchases)
		orders.GET("/sales", h.getUserSales)
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
