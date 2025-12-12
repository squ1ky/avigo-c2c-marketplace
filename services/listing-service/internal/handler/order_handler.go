package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders")
	{
		orders.GET("/purchases", h.getMyPurchases)
		orders.GET("/sales", h.getMySales)
	}
}

func (h *Handler) getMyPurchases(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	resp, err := h.orderSvc.GetUserOrders(c.Request.Context(), userID, false)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getMySales(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	resp, err := h.orderSvc.GetUserOrders(c.Request.Context(), userID, true)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
