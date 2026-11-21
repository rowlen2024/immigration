package handler

import (
	"net/http"
	"strconv"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListPublishedProjectRegions(c *gin.Context) {
	rawIDs := c.Query("ids")
	if rawIDs == "" {
		c.JSON(http.StatusOK, dto.Success([]dto.PublishedProjectRegion{}))
		return
	}
	const maxProjectIDs = 500
	parts := strings.Split(rawIDs, ",")
	if len(parts) > maxProjectIDs {
		c.JSON(http.StatusBadRequest, dto.Error(400, "too many project ids"))
		return
	}
	ids := make([]uint64, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseUint(part, 10, 64)
		if err != nil || id == 0 {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
			return
		}
		ids = append(ids, id)
	}
	items, err := h.svc.Project.PublishedProjectRegions(ids)
	if err != nil {
		logging.Logger.Error("failed in ListPublishedProjectRegions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) AdminListProjectRegions(c *gin.Context) {
	items, err := h.svc.Project.RegionOptions()
	if err != nil {
		logging.Logger.Error("failed in AdminListProjectRegions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}
