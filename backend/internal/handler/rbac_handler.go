package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PermissionResolver() *service.RBACService {
	return h.svc.RBAC
}

func (h *Handler) ListPermissions(c *gin.Context) {
	permissions, err := h.svc.RBAC.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(permissions))
}

func (h *Handler) MyPermissions(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Error(401, "unauthorized"))
		return
	}
	userID, ok := userIDValue.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Error(401, "unauthorized"))
		return
	}
	permissions, err := h.svc.RBAC.EffectivePermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{"permissions": permissions}))
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.svc.RBAC.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(roles))
}

func (h *Handler) GetRole(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid role id"))
		return
	}
	role, err := h.svc.RBAC.GetRole(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "role not found"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(role))
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	role, err := h.svc.RBAC.CreateRole(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(role))
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid role id"))
		return
	}
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	role, err := h.svc.RBAC.UpdateRole(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(role))
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid role id"))
		return
	}
	if err := h.svc.RBAC.DeleteRole(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

func (h *Handler) SaveRolePermissions(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid role id"))
		return
	}
	var req dto.SaveRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	if err := h.svc.RBAC.ReplaceRolePermissions(id, req.PermissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	role, err := h.svc.RBAC.GetRole(id)
	if err != nil {
		c.JSON(http.StatusOK, dto.Success(nil))
		return
	}
	c.JSON(http.StatusOK, dto.Success(role))
}
