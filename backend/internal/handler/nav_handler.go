package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetNavigation(c *gin.Context) {
	position := c.DefaultQuery("position", "header")
	if position != "header" && position != "footer" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "position must be 'header' or 'footer'"))
		return
	}
	tree, err := h.svc.Nav.GetTree(position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(tree))
}

func (h *Handler) AdminListNavigationTree(c *gin.Context) {
	tree, err := h.svc.Nav.GetAdminTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(tree))
}

func (h *Handler) AdminListNavigation(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.svc.Nav.AdminList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, page, pageSize, total))
}

func (h *Handler) CreateNavigation(c *gin.Context) {
	var nav model.Navigation
	if err := c.ShouldBindJSON(&nav); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	created, err := h.svc.Nav.Create(&nav)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateNavigation(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid navigation id"))
		return
	}
	var nav model.Navigation
	if err := c.ShouldBindJSON(&nav); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.Nav.Update(id, &nav)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteNavigation(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid navigation id"))
		return
	}
	if err := h.svc.Nav.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}
