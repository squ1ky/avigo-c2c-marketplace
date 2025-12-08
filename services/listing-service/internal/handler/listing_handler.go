package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
)

func (h *Handler) initListingRoutes(api *gin.RouterGroup) {
	listings := api.Group("/listings")
	{
		listings.POST("", h.createListing)
		listings.GET("/search", h.searchListings)
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

func (h *Handler) searchListings(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		badRequest(c, "INVALID_PAGE", err)
		return
	}
	if page <= 0 {
		badRequest(c, "INVALID_PAGE", fmt.Errorf("page must be positive"))
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		badRequest(c, "INVALID_LIMIT", err)
		return
	}
	if limit <= 0 {
		badRequest(c, "INVALID_LIMIT", fmt.Errorf("limit must be positive"))
		return
	}

	var categoryID *uuid.UUID
	if categoryStr := c.Query("category_id"); categoryStr != "" {
		categoryUUID, err := uuid.Parse(categoryStr)
		if err != nil {
			badRequest(c, "INVALID_CATEGORY_ID", err)
			return
		}
		categoryID = &categoryUUID
	}

	var minPrice, maxPrice *float64
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		val, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			badRequest(c, "INVALID_MIN_PRICE", err)
			return
		}
		minPrice = &val
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		val, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			badRequest(c, "INVALID_MAX_PRICE", err)
			return
		}
		maxPrice = &val
	}

	input := dto.SearchListingsInput{
		Query:      c.Query("q"),
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
	}

	resp, err := h.listingSvc.Search(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
