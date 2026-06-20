package handler

import (
	"net/http"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PublicVersions(c *gin.Context) {
	rawKeys := c.Query("keys")
	if rawKeys == "" {
		c.JSON(http.StatusOK, dto.Success(map[string]string{}))
		return
	}

	keys := strings.Split(rawKeys, ",")
	versions, err := h.svc.PublicVersion.Resolve(keys)
	if err != nil {
		logging.Logger.Error("failed in PublicVersions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(versions))
}
