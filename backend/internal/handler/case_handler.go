package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListCases(c *gin.Context) {
	cases, err := h.svc.Case.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cases))
}

func (h *Handler) AdminListCases(c *gin.Context) {
	page, perPage := parsePagination(c)
	search := c.Query("search")

	cases, total, err := h.svc.Case.AdminList(page, perPage, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(cases, page, perPage, total))
}

func (h *Handler) CreateCase(c *gin.Context) {
	var caseModel model.Case
	if err := c.ShouldBindJSON(&caseModel); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	created, err := h.svc.Case.Create(&caseModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateCase(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	var caseModel model.Case
	if err := c.ShouldBindJSON(&caseModel); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.Case.Update(id, &caseModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteCase(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	if err := h.svc.Case.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
