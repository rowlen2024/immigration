package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	users, total, err := h.svc.User.List(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListUsers", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(users, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(users))
	}
}

func (h *Handler) AdminGetUser(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid user id"))
		return
	}

	user, err := h.svc.User.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "user not found"))
		return
	}

	if h.svc.RBAC == nil {
		c.JSON(http.StatusOK, dto.Success(user))
		return
	}

	overrides, err := h.svc.RBAC.UserOverrides(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"user":                 user,
		"permission_overrides": overrides,
	}))
}

func (h *Handler) AdminCreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	user, err := h.svc.User.Create(req.Username, req.Password, req.DisplayName, req.Role, req.PermissionOverrides)
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

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid user id"))
		return
	}

	if err := h.svc.User.Delete(id); err != nil {
		logging.Logger.Warn("business error in AdminDeleteUser", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
