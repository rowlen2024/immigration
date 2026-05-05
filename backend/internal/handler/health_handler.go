package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Success(gin.H{"status": "ok"}))
}
