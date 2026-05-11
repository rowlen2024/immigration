package router

import (
	"time"

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
		api.GET("/health", h.Health)

		api.POST("/leads", middleware.RateLimit(10, time.Minute), h.CreateLead)

		api.GET("/projects", h.ListProjects)
		api.GET("/projects/:slug", h.GetProject)
		api.GET("/projects/compare", h.CompareProjects)
		api.GET("/faqs", h.ListFAQs)
		api.GET("/pages", h.ListPages)
		api.GET("/pages/*slug", h.GetPage)
		api.GET("/cases", h.ListCases)
		api.GET("/home-config", h.GetHomeConfig)
		api.GET("/navigation", h.GetNavigation)
		api.GET("/search", h.Search)
		api.GET("/site-config", h.GetSiteConfig)

		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimit(5, time.Minute), h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.Auth(cfg))
		{
			admin.GET("/dashboard/stats", middleware.RBAC("admin:read"), h.DashboardStats)

			projects := admin.Group("/projects")
			{
				projects.GET("", middleware.RBAC("admin:read"), h.AdminListProjects)
				projects.POST("", middleware.RBAC("projects:write"), h.CreateProject)
				projects.PUT("/:id", middleware.RBAC("projects:write"), h.UpdateProject)
				projects.DELETE("/:id", middleware.RBAC("projects:write"), h.DeleteProject)

				projects.GET("/:id/requirements", middleware.RBAC("admin:read"), h.ListRequirements)
				projects.POST("/:id/requirements", middleware.RBAC("projects:write"), h.CreateRequirement)
				projects.PUT("/:id/requirements/:rid", middleware.RBAC("projects:write"), h.UpdateRequirement)
				projects.DELETE("/:id/requirements/:rid", middleware.RBAC("projects:write"), h.DeleteRequirement)

				projects.GET("/:id/cost-items", middleware.RBAC("admin:read"), h.ListCostItems)
				projects.POST("/:id/cost-items", middleware.RBAC("projects:write"), h.CreateCostItem)
				projects.PUT("/:id/cost-items/:cid", middleware.RBAC("projects:write"), h.UpdateCostItem)
				projects.DELETE("/:id/cost-items/:cid", middleware.RBAC("projects:write"), h.DeleteCostItem)

				projects.GET("/:id/timeline-phases", middleware.RBAC("admin:read"), h.ListTimelinePhases)
				projects.POST("/:id/timeline-phases", middleware.RBAC("projects:write"), h.CreateTimelinePhase)
				projects.PUT("/:id/timeline-phases/:tid", middleware.RBAC("projects:write"), h.UpdateTimelinePhase)
				projects.DELETE("/:id/timeline-phases/:tid", middleware.RBAC("projects:write"), h.DeleteTimelinePhase)

				// Cases sub-resources
				projects.GET("/:id/cases", middleware.RBAC("admin:read"), h.ListProjectCases)
				projects.POST("/:id/cases", middleware.RBAC("projects:write"), h.CreateProjectCase)
				projects.PUT("/:id/cases/:cid", middleware.RBAC("projects:write"), h.UpdateProjectCase)
				projects.DELETE("/:id/cases/:cid", middleware.RBAC("projects:write"), h.DeleteProjectCase)

				// News sub-resources
				projects.GET("/:id/news", middleware.RBAC("admin:read"), h.ListProjectNews)
				projects.POST("/:id/news", middleware.RBAC("projects:write"), h.AddProjectNews)
				projects.DELETE("/:id/news/:page_id", middleware.RBAC("projects:write"), h.RemoveProjectNews)

				// Compare config
				projects.GET("/:id/compare-config", middleware.RBAC("admin:read"), h.GetCompareConfig)
				projects.PUT("/:id/compare-config", middleware.RBAC("projects:write"), h.SaveCompareConfig)
				}

			admin.GET("/faqs", middleware.RBAC("admin:read"), h.AdminListFAQs)
			admin.POST("/faqs", middleware.RBAC("content:write"), h.CreateFAQ)
			admin.PUT("/faqs/:id", middleware.RBAC("content:write"), h.UpdateFAQ)
			admin.DELETE("/faqs/:id", middleware.RBAC("content:write"), h.DeleteFAQ)

			admin.GET("/pages", middleware.RBAC("admin:read"), h.AdminListPages)
			admin.POST("/pages", middleware.RBAC("content:write"), h.CreatePage)
			admin.PUT("/pages/:id", middleware.RBAC("content:write"), h.UpdatePage)
			admin.DELETE("/pages/:id", middleware.RBAC("content:write"), h.DeletePage)

			admin.GET("/cases", middleware.RBAC("admin:read"), h.AdminListCases)
			admin.POST("/cases", middleware.RBAC("content:write"), h.CreateCase)
			admin.PUT("/cases/:id", middleware.RBAC("content:write"), h.UpdateCase)
			admin.DELETE("/cases/:id", middleware.RBAC("content:write"), h.DeleteCase)

			admin.GET("/compare-fields", middleware.RBAC("admin:read"), h.ListCompareFields)
			admin.GET("/leads", middleware.RBAC("leads:read"), h.AdminListLeads)
			admin.PUT("/leads/:id", middleware.RBAC("leads:read"), h.UpdateLead)

			admin.GET("/users", middleware.RBAC("admin:write"), h.AdminListUsers)
			admin.POST("/users", middleware.RBAC("admin:write"), h.AdminCreateUser)
			admin.PUT("/users/:id", middleware.RBAC("admin:write"), h.AdminUpdateUser)

			admin.POST("/media/upload", middleware.RBAC("content:write"), h.UploadMedia)
			admin.GET("/media", middleware.RBAC("content:write"), h.ListMedia)
			admin.DELETE("/media/:id", middleware.RBAC("content:write"), h.DeleteMedia)

			admin.GET("/home-config", middleware.RBAC("admin:read"), h.GetHomeConfig)
			admin.PUT("/home-config", middleware.RBAC("content:write"), h.UpdateHomeConfig)

			admin.PUT("/site-config", middleware.RBAC("content:write"), h.UpdateSiteConfig)

			admin.GET("/navigation", middleware.RBAC("admin:read"), h.AdminListNavigationTree)
			admin.POST("/navigation", middleware.RBAC("content:write"), h.CreateNavigation)
			admin.PUT("/navigation/:id", middleware.RBAC("content:write"), h.UpdateNavigation)
			admin.DELETE("/navigation/:id", middleware.RBAC("content:write"), h.DeleteNavigation)
		}
	}

	return r
}
