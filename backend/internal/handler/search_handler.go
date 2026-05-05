package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"

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
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(results))
}
