package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

// Requirements

func (h *Handler) ListRequirements(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.Requirement.List(projectID)
	if err != nil {
		logging.Logger.Error("failed in ListRequirements", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateRequirement(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.Requirement
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	item.ID = 0 // ensure DB auto-increment
	created, err := h.svc.Requirement.Create(projectID, &item)
	if err != nil {
		logging.Logger.Error("failed in CreateRequirement", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateRequirement(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	rid, err := parseIDParam(c, "rid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid requirement id"))
		return
	}
	var item model.Requirement
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.Requirement.Update(projectID, rid, &item)
	if err != nil {
		logging.Logger.Error("failed in UpdateRequirement", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteRequirement(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	rid, err := parseIDParam(c, "rid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid requirement id"))
		return
	}
	if err := h.svc.Requirement.Delete(projectID, rid); err != nil {
		logging.Logger.Error("failed in DeleteRequirement", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

// Cost Items

func (h *Handler) ListCostItems(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.CostItem.List(projectID)
	if err != nil {
		logging.Logger.Error("failed in ListCostItems", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateCostItem(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.CostItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	item.ID = 0 // ensure DB auto-increment
	created, err := h.svc.CostItem.Create(projectID, &item)
	if err != nil {
		logging.Logger.Error("failed in CreateCostItem", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateCostItem(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	cid, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid cost item id"))
		return
	}
	var item model.CostItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.CostItem.Update(projectID, cid, &item)
	if err != nil {
		logging.Logger.Error("failed in UpdateCostItem", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteCostItem(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	cid, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid cost item id"))
		return
	}
	if err := h.svc.CostItem.Delete(projectID, cid); err != nil {
		logging.Logger.Error("failed in DeleteCostItem", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

// Timeline Phases

func (h *Handler) ListTimelinePhases(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.TimelinePhase.List(projectID)
	if err != nil {
		logging.Logger.Error("failed in ListTimelinePhases", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateTimelinePhase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.TimelinePhase
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	item.ID = 0 // ensure DB auto-increment
	created, err := h.svc.TimelinePhase.Create(projectID, &item)
	if err != nil {
		logging.Logger.Error("failed in CreateTimelinePhase", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateTimelinePhase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	tid, err := parseIDParam(c, "tid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid timeline phase id"))
		return
	}
	var item model.TimelinePhase
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.TimelinePhase.Update(projectID, tid, &item)
	if err != nil {
		logging.Logger.Error("failed in UpdateTimelinePhase", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteTimelinePhase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	tid, err := parseIDParam(c, "tid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid timeline phase id"))
		return
	}
	if err := h.svc.TimelinePhase.Delete(projectID, tid); err != nil {
		logging.Logger.Error("failed in DeleteTimelinePhase", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}
	// Project Advantages

	func (h *Handler) ListProjectAdvantages(c *gin.Context) {
		projectID, err := parseIDParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
			return
		}
		items, err := h.svc.Advantage.List(projectID)
		if err != nil {
			logging.Logger.Error("failed in DeleteTimelinePhase", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(items))
	}

	func (h *Handler) CreateProjectAdvantage(c *gin.Context) {
		projectID, err := parseIDParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
			return
		}
		var item model.ProjectAdvantage
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
			return
		}
		item.ID = 0 // ensure DB auto-increment
		created, err := h.svc.Advantage.Create(projectID, &item)
		if err != nil {
			logging.Logger.Error("failed in DeleteTimelinePhase", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusCreated, dto.Success(created))
	}

	func (h *Handler) UpdateProjectAdvantage(c *gin.Context) {
		projectID, err := parseIDParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
			return
		}
		aid, err := parseIDParam(c, "aid")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid advantage id"))
			return
		}
		var item model.ProjectAdvantage
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
			return
		}
		updated, err := h.svc.Advantage.Update(projectID, aid, &item)
		if err != nil {
			logging.Logger.Error("failed in DeleteTimelinePhase", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(updated))
	}

	func (h *Handler) DeleteProjectAdvantage(c *gin.Context) {
		projectID, err := parseIDParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
			return
		}
		aid, err := parseIDParam(c, "aid")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid advantage id"))
			return
		}
		if err := h.svc.Advantage.Delete(projectID, aid); err != nil {
			logging.Logger.Error("failed in DeleteTimelinePhase", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(nil))
	}
