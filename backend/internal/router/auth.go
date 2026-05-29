package router

import (
	"time"

	"mygo-immigration/backend/internal/handler"
	"mygo-immigration/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(api *gin.RouterGroup, h *handler.Handler) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", middleware.RateLimit(5, time.Minute), h.Login)
		auth.POST("/refresh", h.RefreshToken)
	}
}
