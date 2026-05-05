package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	tokenPair, err := h.svc.Auth.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Error(401, err.Error()))
		return
	}

	c.SetCookie(
		"refresh_token", tokenPair.RefreshToken,
		604800, // 7 days
		"/api/v1/auth",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, dto.Success(loginResponse{AccessToken: tokenPair.AccessToken}))
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Error(401, "refresh token required"))
		return
	}

	tokenPair, err := h.svc.Auth.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Error(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(loginResponse{AccessToken: tokenPair.AccessToken}))
}
