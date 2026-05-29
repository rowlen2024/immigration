package handler

import (
	"net/http"
	"strconv"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListFAQs(c *gin.Context) {
	page, perPage := parsePagination(c)

	var projectID *uint64
	if v := c.Query("project_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			projectID = &id
		}
	}

	var isGlobal *bool
	if v := c.Query("is_global"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			isGlobal = &b
		}
	}

	faqs, total, err := h.svc.FAQ.List(projectID, isGlobal, page, perPage)
	if err != nil {
		logging.Logger.Error("failed in ListFAQs", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, page, perPage, total))
}

func (h *Handler) AdminListFAQs(c *gin.Context) {
	var projectID *uint64
	if v := c.Query("project_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			projectID = &id
		}
	}

	search := c.Query("search")

	if c.Query("all") == "true" {
		faqs, err := h.svc.FAQ.ListAll(projectID, search)
		if err != nil {
			logging.Logger.Error("failed in AdminListFAQs", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(faqs))
		return
	}

	page, perPage := parsePagination(c)
	faqs, total, err := h.svc.FAQ.AdminList(projectID, search, page, perPage)
	if err != nil {
		logging.Logger.Error("failed in AdminListFAQs", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, page, perPage, total))
}

func (h *Handler) ListFAQProjects(c *gin.Context) {
	projects, err := h.svc.FAQ.ListProjects()
	if err != nil {
		logging.Logger.Error("failed in ListFAQProjects", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(projects))
}

func (h *Handler) CreateFAQ(c *gin.Context) {
	var faq model.FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	faq.ID = 0 // ensure DB auto-increment
	created, err := h.svc.FAQ.Create(&faq)
	if err != nil {
		logging.Logger.Warn("business error in CreateFAQ", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateFAQ(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid faq id"))
		return
	}

	var req dto.UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.FAQ.Update(id, req)
	if err != nil {
		logging.Logger.Warn("business error in UpdateFAQ", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteFAQ(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid faq id"))
		return
	}

	if err := h.svc.FAQ.Delete(id); err != nil {
		logging.Logger.Warn("business error in DeleteFAQ", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
