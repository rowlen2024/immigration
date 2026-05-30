package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListFAQs(c *gin.Context) {
	var req dto.FAQListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	faqs, total, err := h.svc.FAQ.List(req)
	if err != nil {
		logging.Logger.Error("failed in ListFAQs", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(faqs))
	}
}

func (h *Handler) AdminListFAQs(c *gin.Context) {
	var req dto.FAQListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	faqs, total, err := h.svc.FAQ.List(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListFAQs", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(faqs))
	}
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

	faq.ID = 0
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
