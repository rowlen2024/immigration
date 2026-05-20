package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

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
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	user, err := h.svc.User.Create(req.Username, req.Password, req.DisplayName, req.Role)
	if err != nil {
		logging.Logger.Error("failed in AdminCreateUser", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
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

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	user, err := h.svc.User.Update(id, updates)
	if err != nil {
		logging.Logger.Error("failed in AdminUpdateUser", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(user))
}
