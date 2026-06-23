package router

import (
	"mygo-immigration/backend/internal/handler"
	"mygo-immigration/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdminProjectRoutes(admin *gin.RouterGroup, h *handler.Handler) {
	projects := admin.Group("/projects")
	{
		projects.GET("", middleware.RBAC("projects:read"), h.AdminListProjects)
		projects.GET("/options", middleware.RBAC("projects:read"), h.AdminListProjectOptions)
		projects.POST("", middleware.RBAC("projects:write"), h.CreateProject)
		projects.PUT("/:id", middleware.RBAC("projects:write"), h.UpdateProject)
		projects.DELETE("/:id", middleware.RBAC("projects:write"), h.DeleteProject)

		projects.GET("/:id/requirements", middleware.RBAC("projects:read"), h.ListRequirements)
		projects.POST("/:id/requirements", middleware.RBAC("projects:write"), h.CreateRequirement)
		projects.PUT("/:id/requirements/:rid", middleware.RBAC("projects:write"), h.UpdateRequirement)
		projects.DELETE("/:id/requirements/:rid", middleware.RBAC("projects:write"), h.DeleteRequirement)

		projects.GET("/:id/cost-items", middleware.RBAC("projects:read"), h.ListCostItems)
		projects.POST("/:id/cost-items", middleware.RBAC("projects:write"), h.CreateCostItem)
		projects.PUT("/:id/cost-items/:cid", middleware.RBAC("projects:write"), h.UpdateCostItem)
		projects.DELETE("/:id/cost-items/:cid", middleware.RBAC("projects:write"), h.DeleteCostItem)

		projects.GET("/:id/timeline-phases", middleware.RBAC("projects:read"), h.ListTimelinePhases)
		projects.POST("/:id/timeline-phases", middleware.RBAC("projects:write"), h.CreateTimelinePhase)
		projects.PUT("/:id/timeline-phases/:tid", middleware.RBAC("projects:write"), h.UpdateTimelinePhase)
		projects.DELETE("/:id/timeline-phases/:tid", middleware.RBAC("projects:write"), h.DeleteTimelinePhase)

		projects.GET("/:id/advantages", middleware.RBAC("projects:read"), h.ListProjectAdvantages)
		projects.POST("/:id/advantages", middleware.RBAC("projects:write"), h.CreateProjectAdvantage)
		projects.PUT("/:id/advantages/:aid", middleware.RBAC("projects:write"), h.UpdateProjectAdvantage)
		projects.DELETE("/:id/advantages/:aid", middleware.RBAC("projects:write"), h.DeleteProjectAdvantage)

		projects.GET("/:id/cases", middleware.RBAC("cases:read"), h.ListProjectCases)
		projects.POST("/:id/cases", middleware.RBAC("cases:write"), h.CreateProjectCase)
		projects.PUT("/:id/cases/:cid", middleware.RBAC("cases:write"), h.UpdateProjectCase)
		projects.DELETE("/:id/cases/:cid", middleware.RBAC("cases:write"), h.DeleteProjectCase)

		projects.GET("/:id/testimonials", middleware.RBAC("testimonials:read"), h.ListProjectTestimonials)
		projects.POST("/:id/testimonials", middleware.RBAC("testimonials:write"), h.CreateProjectTestimonial)
		projects.PUT("/:id/testimonials/:tid", middleware.RBAC("testimonials:write"), h.UpdateProjectTestimonial)
		projects.DELETE("/:id/testimonials/:tid", middleware.RBAC("testimonials:write"), h.DeleteProjectTestimonial)

		projects.GET("/:id/news", middleware.RBAC("pages:read"), h.ListProjectNews)
		projects.POST("/:id/news", middleware.RBAC("pages:write"), h.AddProjectNews)
		projects.DELETE("/:id/news/:page_id", middleware.RBAC("pages:write"), h.RemoveProjectNews)

		projects.GET("/:id/compare-config", middleware.RBAC("projects:read"), h.GetCompareConfig)
		projects.PUT("/:id/compare-config", middleware.RBAC("projects:write"), h.SaveCompareConfig)
	}
}
