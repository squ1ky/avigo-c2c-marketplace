package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/service"
)

type Handler struct {
	listingSvc *service.ListingService
	orderSvc   *service.OrderService
	reviewSvc  *service.ReviewService
	mediaSvc   *service.MediaService
}

func NewHandler(
	listingSve *service.ListingService,
	orderSvc *service.OrderService,
	reviewSvc *service.ReviewService,
	mediaSvc *service.MediaService,
) *Handler {
	return &Handler{
		listingSvc: listingSve,
		orderSvc:   orderSvc,
		reviewSvc:  reviewSvc,
		mediaSvc:   mediaSvc,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initListingRoutes(v1)
		h.initOrderRoutes(v1)
		//h.initReviewRoutes(v1)
		h.initMediaRoutes(v1)
	}
}
