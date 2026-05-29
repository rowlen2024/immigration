package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

type updateLeadRequest struct {
	Status string  `json:"status" binding:"required"`
	Notes  *string `json:"notes"`
}

func (h *Handler) CreateLead(c *gin.Context) {
	var req dto.LeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	lead, err := h.svc.Lead.Create(&req)
	if err != nil {
		logging.Logger.Warn("business error in CreateLead", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(lead))
}

func (h *Handler) AdminListLeads(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = defaultPage
		pagination.PerPage = defaultPerPage
	}

	page := pagination.Page
	if page < 1 {
		page = defaultPage
	}

	perPage := pagination.PerPage
	if perPage < 1 {
		perPage = defaultPerPage
	}

	leads, total, err := h.svc.Lead.AdminList(page, perPage, pagination.Status)
	if err != nil {
		logging.Logger.Error("failed in AdminListLeads", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(leads, page, perPage, total))
}

func (h *Handler) UpdateLead(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid lead id"))
		return
	}

	var req updateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	lead, err := h.svc.Lead.Update(id, req.Status, notes)
	if err != nil {
		logging.Logger.Warn("business error in UpdateLead", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(lead))
}
