package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initMediaRoutes(api *gin.RouterGroup) {
	media := api.Group("/media")
	{
		media.POST("/upload", h.uploadMedia)
	}
}

func (h *Handler) uploadMedia(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	// Limit 10MB
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		badRequest(c, "FILE_TOO_LARGE", err)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		badRequest(c, "MISSING_FILE", err)
		return
	}
	defer file.Close()

	res, err := h.mediaSvc.UploadTempFile(c.Request.Context(), userID, header)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
