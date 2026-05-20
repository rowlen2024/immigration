package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "query parameter q is required"))
		return
	}

	results, err := h.svc.Search.Search(q)
	if err != nil {
		logging.Logger.Error("failed in Search", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(results))
}
