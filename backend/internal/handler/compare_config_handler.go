package handler

import (
	"encoding/json"
	"net/http"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// GetCompareConfig returns the compare config for a project.
func (h *Handler) GetCompareConfig(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	cfg, err := h.svc.CompareConfig.GetByProjectID(projectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.Success(nil))
		return
	}
	c.JSON(http.StatusOK, dto.Success(cfg))
}

type saveCompareConfigRequest struct {
	ProjectID     uint64   `json:"project_id"`
	CompareWith   []string `json:"compare_with"`
	CompareFields []string `json:"compare_fields"`
}

// SaveCompareConfig upserts the compare config for a project.
// If compare_with has fewer than 2 items, the config is deleted.
func (h *Handler) SaveCompareConfig(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var req saveCompareConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	if len(req.CompareWith) < 2 {
		_ = h.svc.CompareConfig.DeleteByProjectID(projectID)
		c.JSON(http.StatusOK, dto.Success(nil))
		return
	}

	cfg := &model.CompareConfig{
		ProjectID:     projectID,
		CompareWith:   marshalJSON(req.CompareWith),
		CompareFields: marshalJSON(req.CompareFields),
	}

	saved, err := h.svc.CompareConfig.Save(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(saved))
}

// ListCompareFields returns available compare field definitions.
func (h *Handler) ListCompareFields(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Success(config.CompareFields))
}

func marshalJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}
