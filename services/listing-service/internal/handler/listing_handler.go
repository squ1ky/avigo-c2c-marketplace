package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"net/http"
)

func (h *Handler) initListingRoutes(api *gin.RouterGroup) {
	listings := api.Group("/listings")
	{
		listings.POST("", h.createListing)
		listings.GET("/:id", h.getListing)
		listings.PUT("/:id", h.updateListing)
		listings.DELETE("/:id", h.deleteListing)
	}
}

func (h *Handler) createListing(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	var input dto.CreateListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		invalidJSON(c, err)
		return
	}

	input.UserID = userID

	resp, err := h.listingSvc.Create(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) getListing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		badRequest(c, "INVALID_ID", err)
		return
	}

	resp, err := h.listingSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) updateListing(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	listingID, err := uuid.Parse(idStr)
	if err != nil {
		badRequest(c, "INVALID_ID", err)
		return
	}

	var input dto.UpdateListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		invalidJSON(c, err)
		return
	}

	input.UserID = userID
	input.ID = listingID

	resp, err := h.listingSvc.Update(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) deleteListing(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		badRequest(c, "INVALID_ID", err)
		return
	}

	if err := h.listingSvc.Delete(c.Request.Context(), id, userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
