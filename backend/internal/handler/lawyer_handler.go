package handler

import (
	"net/http"
	"strconv"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/lawyers — public
func (h *Handler) ListLawyers(c *gin.Context) {
	items, err := h.svc.Lawyer.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

// GET /api/v1/admin/lawyers — admin list
func (h *Handler) AdminListLawyers(c *gin.Context) {
	items, err := h.svc.Lawyer.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

// GET /api/v1/admin/lawyers/:id
func (h *Handler) AdminGetLawyer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid id"))
		return
	}
	item, err := h.svc.Lawyer.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(item))
}

// POST /api/v1/admin/lawyers
func (h *Handler) CreateLawyer(c *gin.Context) {
	var input service.CreateLawyerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request body"))
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "name is required"))
		return
	}
	item, err := h.svc.Lawyer.Create(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(item))
}

// PUT /api/v1/admin/lawyers/:id
func (h *Handler) UpdateLawyer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid id"))
		return
	}
	var input service.CreateLawyerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request body"))
		return
	}
	item, err := h.svc.Lawyer.Update(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(item))
}

// DELETE /api/v1/admin/lawyers/:id
func (h *Handler) DeleteLawyer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid id"))
		return
	}
	if err := h.svc.Lawyer.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}
