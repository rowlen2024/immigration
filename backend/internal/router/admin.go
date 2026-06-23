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
	admin.Use(middleware.LoadPermissions(h.PermissionResolver()))
	{
		admin.GET("/dashboard/stats", middleware.RBAC("dashboard:read"), h.DashboardStats)

		registerAdminProjectRoutes(admin, h)

		admin.GET("/faqs", middleware.RBAC("faqs:read"), h.AdminListFAQs)
		admin.POST("/faqs", middleware.RBAC("faqs:write"), h.CreateFAQ)
		admin.PUT("/faqs/:id", middleware.RBAC("faqs:write"), h.UpdateFAQ)
		admin.DELETE("/faqs/:id", middleware.RBAC("faqs:write"), h.DeleteFAQ)

		admin.GET("/pages/preview", middleware.RBAC("pages:read"), h.PreviewPage)
		admin.GET("/pages/options", middleware.RBAC("pages:read"), h.AdminListPageOptions)
		admin.GET("/pages", middleware.RBAC("pages:read"), h.AdminListPages)
		admin.POST("/pages", middleware.RBAC("pages:write"), h.CreatePage)
		admin.PUT("/pages/:id", middleware.RBAC("pages:write"), h.UpdatePage)
		admin.DELETE("/pages/:id", middleware.RBAC("pages:write"), h.DeletePage)

		admin.GET("/cases/options", middleware.RBAC("cases:read"), h.AdminListCaseOptions)
		admin.GET("/cases", middleware.RBAC("cases:read"), h.AdminListCases)
		admin.POST("/cases", middleware.RBAC("cases:write"), h.CreateCase)
		admin.PUT("/cases/:id", middleware.RBAC("cases:write"), h.UpdateCase)
		admin.DELETE("/cases/:id", middleware.RBAC("cases:write"), h.DeleteCase)

		admin.GET("/lawyers", middleware.RBAC("lawyers:read"), h.AdminListLawyers)
		admin.GET("/lawyers/:id", middleware.RBAC("lawyers:read"), h.AdminGetLawyer)
		admin.POST("/lawyers", middleware.RBAC("lawyers:write"), h.CreateLawyer)
		admin.PUT("/lawyers/:id", middleware.RBAC("lawyers:write"), h.UpdateLawyer)
		admin.DELETE("/lawyers/:id", middleware.RBAC("lawyers:write"), h.DeleteLawyer)

		admin.GET("/testimonials/options", middleware.RBAC("testimonials:read"), h.AdminListTestimonialOptions)
		admin.GET("/testimonials", middleware.RBAC("testimonials:read"), h.AdminListTestimonials)

		admin.GET("/compare-fields", middleware.RBAC("projects:read"), h.ListCompareFields)
		admin.GET("/leads", middleware.RBAC("leads:read"), h.AdminListLeads)
		admin.PUT("/leads/:id", middleware.RBAC("leads:write"), h.UpdateLead)

		admin.GET("/users", middleware.RBAC("users:read"), h.AdminListUsers)
		admin.GET("/users/:id", middleware.RBAC("users:read"), h.AdminGetUser)
		admin.POST("/users", middleware.RBAC("users:write"), h.AdminCreateUser)
		admin.PUT("/users/:id", middleware.RBAC("users:write"), h.AdminUpdateUser)
		admin.DELETE("/users/:id", middleware.RBAC("users:write"), h.AdminDeleteUser)

		admin.POST("/media/upload", middleware.RBAC("media:write"), h.UploadMedia)
		admin.GET("/media", middleware.RBAC("media:read"), h.ListMedia)
		admin.DELETE("/media/:id", middleware.RBAC("media:write"), h.DeleteMedia)
		admin.GET("/media/unused", middleware.RBAC("media:read"), h.FindUnusedMedia)
		admin.POST("/media/cleanup", middleware.RBAC("media:write"), h.CleanupUnusedMedia)

		admin.GET("/home-config", middleware.RBAC("homepage:read"), h.GetAdminHomeConfig)
		admin.PUT("/home-config", middleware.RBAC("homepage:write"), h.UpdateHomeConfig)

		admin.GET("/site-config", middleware.RBAC("settings:read"), h.GetSiteConfig)
		admin.PUT("/site-config", middleware.RBAC("settings:write"), h.UpdateSiteConfig)

		admin.GET("/navigation", middleware.RBAC("navigation:read"), h.AdminListNavigationTree)
		admin.POST("/navigation", middleware.RBAC("navigation:write"), h.CreateNavigation)
		admin.PUT("/navigation/:id", middleware.RBAC("navigation:write"), h.UpdateNavigation)
		admin.DELETE("/navigation/:id", middleware.RBAC("navigation:write"), h.DeleteNavigation)

		admin.GET("/permissions", middleware.RBACAny("roles:read", "users:read"), h.ListPermissions)
		admin.GET("/me/permissions", h.MyPermissions)
		admin.GET("/roles", middleware.RBACAny("roles:read", "users:read"), h.ListRoles)
		admin.GET("/roles/:id", middleware.RBAC("roles:read"), h.GetRole)
		admin.POST("/roles", middleware.RBAC("roles:write"), h.CreateRole)
		admin.PUT("/roles/:id", middleware.RBAC("roles:write"), h.UpdateRole)
		admin.DELETE("/roles/:id", middleware.RBAC("roles:write"), h.DeleteRole)
		admin.PUT("/roles/:id/permissions", middleware.RBAC("roles:write"), h.SaveRolePermissions)
	}
}
