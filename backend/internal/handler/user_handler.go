package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminListUsers(c *gin.Context) {
	page, perPage := parsePagination(c)

	users, total, err := h.svc.User.ListPaginated(page, perPage)
	if err != nil {
		logging.Logger.Error("failed in AdminListUsers", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(users, page, perPage, total))
}

func (h *Handler) AdminCreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	user, err := h.svc.User.Create(req.Username, req.Password, req.DisplayName, req.Role)
	if err != nil {
		logging.Logger.Warn("business error in AdminCreateUser", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(user))
}

func (h *Handler) AdminUpdateUser(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid user id"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	user, err := h.svc.User.Update(id, req)
	if err != nil {
		logging.Logger.Warn("business error in AdminUpdateUser", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(user))
}
