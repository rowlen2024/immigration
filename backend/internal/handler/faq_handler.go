package handler

import (
	"net/http"
	"strconv"

	"mygo-immigration/backend/internal/dto"
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
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, page, perPage, total))
}

func (h *Handler) AdminListFAQs(c *gin.Context) {
	page, perPage := parsePagination(c)

	var projectID *uint64
	if v := c.Query("project_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			projectID = &id
		}
	}

	search := c.Query("search")

	faqs, total, err := h.svc.FAQ.AdminList(projectID, search, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, page, perPage, total))
}

func (h *Handler) CreateFAQ(c *gin.Context) {
	var faq model.FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	created, err := h.svc.FAQ.Create(&faq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
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

	var faq model.FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.FAQ.Update(id, &faq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
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
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
