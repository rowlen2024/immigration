package router

import (
	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/handler"
	"mygo-immigration/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 50 << 20 // 50MB

	h := handler.New(db, cfg)

	r.Use(middleware.CORS(cfg.CORSOrigin))
	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		registerPublicRoutes(api, h)
		registerAuthRoutes(api, h)
		registerAdminRoutes(api, h, cfg)
	}

	return r
}
