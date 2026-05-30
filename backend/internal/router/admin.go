package router

import (
	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/handler"
	"mygo-immigration/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(api *gin.RouterGroup, h *handler.Handler, cfg *config.Config) {
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(cfg))
	{
		admin.GET("/dashboard/stats", middleware.RBAC("admin:read"), h.DashboardStats)

		registerAdminProjectRoutes(admin, h)

		admin.GET("/faqs", middleware.RBAC("admin:read"), h.AdminListFAQs)
		admin.POST("/faqs", middleware.RBAC("content:write"), h.CreateFAQ)
		admin.PUT("/faqs/:id", middleware.RBAC("content:write"), h.UpdateFAQ)
		admin.DELETE("/faqs/:id", middleware.RBAC("content:write"), h.DeleteFAQ)

		admin.GET("/pages", middleware.RBAC("admin:read"), h.AdminListPages)
		admin.GET("/pages/preview", middleware.RBAC("admin:read"), h.PreviewPage)
		admin.POST("/pages", middleware.RBAC("content:write"), h.CreatePage)
		admin.PUT("/pages/:id", middleware.RBAC("content:write"), h.UpdatePage)
		admin.DELETE("/pages/:id", middleware.RBAC("content:write"), h.DeletePage)

		admin.GET("/cases", middleware.RBAC("admin:read"), h.AdminListCases)
		admin.POST("/cases", middleware.RBAC("content:write"), h.CreateCase)
		admin.PUT("/cases/:id", middleware.RBAC("content:write"), h.UpdateCase)
		admin.DELETE("/cases/:id", middleware.RBAC("content:write"), h.DeleteCase)

		admin.GET("/lawyers", middleware.RBAC("admin:read"), h.AdminListLawyers)
		admin.GET("/lawyers/:id", middleware.RBAC("admin:read"), h.AdminGetLawyer)
		admin.POST("/lawyers", middleware.RBAC("content:write"), h.CreateLawyer)
		admin.PUT("/lawyers/:id", middleware.RBAC("content:write"), h.UpdateLawyer)
		admin.DELETE("/lawyers/:id", middleware.RBAC("content:write"), h.DeleteLawyer)

		admin.GET("/testimonials", middleware.RBAC("admin:read"), h.AdminListTestimonials)

		admin.GET("/compare-fields", middleware.RBAC("admin:read"), h.ListCompareFields)
		admin.GET("/leads", middleware.RBAC("leads:read"), h.AdminListLeads)
		admin.PUT("/leads/:id", middleware.RBAC("leads:read"), h.UpdateLead)

		admin.GET("/users", middleware.RBAC("admin:write"), h.AdminListUsers)
		admin.GET("/users/:id", middleware.RBAC("admin:write"), h.AdminGetUser)
		admin.POST("/users", middleware.RBAC("admin:write"), h.AdminCreateUser)
		admin.PUT("/users/:id", middleware.RBAC("admin:write"), h.AdminUpdateUser)
		admin.DELETE("/users/:id", middleware.RBAC("admin:write"), h.AdminDeleteUser)

		admin.POST("/media/upload", middleware.RBAC("content:write"), h.UploadMedia)
		admin.GET("/media", middleware.RBAC("content:write"), h.ListMedia)
		admin.DELETE("/media/:id", middleware.RBAC("content:write"), h.DeleteMedia)
		admin.GET("/media/unused", middleware.RBAC("content:write"), h.FindUnusedMedia)
		admin.POST("/media/cleanup", middleware.RBAC("content:write"), h.CleanupUnusedMedia)

		admin.GET("/home-config", middleware.RBAC("admin:read"), h.GetAdminHomeConfig)
		admin.PUT("/home-config", middleware.RBAC("content:write"), h.UpdateHomeConfig)

		admin.GET("/site-config", middleware.RBAC("admin:read"), h.GetSiteConfig)
		admin.PUT("/site-config", middleware.RBAC("content:write"), h.UpdateSiteConfig)

		admin.GET("/navigation", middleware.RBAC("admin:read"), h.AdminListNavigationTree)
		admin.POST("/navigation", middleware.RBAC("content:write"), h.CreateNavigation)
		admin.PUT("/navigation/:id", middleware.RBAC("content:write"), h.UpdateNavigation)
		admin.DELETE("/navigation/:id", middleware.RBAC("content:write"), h.DeleteNavigation)
	}
}
