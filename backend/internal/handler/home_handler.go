package handler

import (
	"encoding/json"
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHomeConfig(c *gin.Context) {
	configs, err := h.svc.HomeConfig.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(configs))
}

func (h *Handler) DashboardStats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(stats))
}

func (h *Handler) GetSiteConfig(c *gin.Context) {
	cfg, err := h.svc.HomeConfig.GetSiteConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cfg))
}

func (h *Handler) UpdateSiteConfig(c *gin.Context) {
	var cfg service.SiteConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request body"))
		return
	}

	if err := h.svc.HomeConfig.UpdateSiteConfig(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

func (h *Handler) UpdateHomeConfig(c *gin.Context) {
	var configs map[string]json.RawMessage
	if err := c.ShouldBindJSON(&configs); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	if err := h.svc.HomeConfig.Update(configs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
