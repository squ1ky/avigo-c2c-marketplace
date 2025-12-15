package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"net/http"
	"strconv"
)

func (h *Handler) initListingRoutes(api *gin.RouterGroup) {
	listings := api.Group("/account/:user_id/listings")
	{
		listings.POST("", h.createListing)
		listings.GET("", h.getUserListings)
		listings.GET("/:listing_id", h.getListing)
		listings.PUT("/:listing_id", h.updateListing)
		listings.DELETE("/:listing_id", h.deleteListing)
	}

	categories := api.Group("/categories")
	{
		categories.GET("", h.getCategories)
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
	idStr := c.Param("listing_id")
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

func (h *Handler) getUserListings(c *gin.Context) {
	targetUserID, err := h.getTargetUserID(c)
	if err != nil {
		badRequest(c, "INVALID_USER_ID", err)
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			if val <= limit {
				limit = val
			}
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val > 0 {
			offset = val
		}
	}

	resp, err := h.listingSvc.GetUserListings(c.Request.Context(), targetUserID, limit, offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getCategories(c *gin.Context) {
	categories, err := h.listingSvc.GetCategoriesTree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handler) updateListing(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	idStr := c.Param("listing_id")
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

	idStr := c.Param("listing_id")
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
