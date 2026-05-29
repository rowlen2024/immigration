package router

import (
	"time"

	"mygo-immigration/backend/internal/handler"
	"mygo-immigration/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerPublicRoutes(api *gin.RouterGroup, h *handler.Handler) {
	api.GET("/health", h.Health)

	api.POST("/leads", middleware.RateLimit(10, time.Minute), h.CreateLead)

	api.GET("/projects", h.ListProjects)
	api.GET("/projects/compare", h.CompareProjects)
	api.GET("/projects/:slug", h.GetProject)
	api.GET("/faqs/projects", h.ListFAQProjects)
	api.GET("/faqs", h.ListFAQs)
	api.GET("/pages", h.ListPages)
	api.GET("/pages/*slug", h.GetPage)
	api.GET("/cases/:slug", h.GetCase)
	api.GET("/cases", h.ListCases)
	api.GET("/testimonials", h.ListAllTestimonials)
	api.GET("/lawyers", h.ListLawyers)
	api.GET("/home-config", h.GetHomeConfig)
	api.GET("/navigation", h.GetNavigation)
	api.GET("/search", h.Search)
	api.GET("/site-config", h.GetSiteConfig)
}
