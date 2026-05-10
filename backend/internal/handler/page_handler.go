package handler

import (
	"net/http"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListPages(c *gin.Context) {
	pages, err := h.svc.Page.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(pages))
}

func (h *Handler) GetPage(c *gin.Context) {
	slug := strings.TrimPrefix(c.Param("slug"), "/")
	if slug == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "slug is required"))
		return
	}

	page, err := h.svc.Page.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "page not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(page))
}

func (h *Handler) AdminListPages(c *gin.Context) {
	if c.Query("all") == "true" {
		pages, err := h.svc.Page.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
			return
		}
		c.JSON(http.StatusOK, dto.Success(pages))
		return
	}

	paginationPage, perPage := parsePagination(c)

	pages, total, err := h.svc.Page.AdminList(paginationPage, perPage, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(pages, paginationPage, perPage, total))
}

func (h *Handler) CreatePage(c *gin.Context) {
	var page model.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	created, err := h.svc.Page.Create(&page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdatePage(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid page id"))
		return
	}

	var page model.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.Page.Update(id, &page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeletePage(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid page id"))
		return
	}

	if err := h.svc.Page.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
