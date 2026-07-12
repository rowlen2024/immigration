package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mygo-immigration/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func TestRegisterPublicRoutesRelatedPagesDoesNotConflictWithPageCatchAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	var matchedPath string
	r.Use(func(c *gin.Context) {
		matchedPath = c.FullPath()
		c.Status(http.StatusNoContent)
		c.Abort()
	})

	api := r.Group("/api/v1")
	registerPublicRoutes(api, &handler.Handler{})

	tests := []struct {
		url      string
		fullPath string
	}{
		{url: "/api/v1/related-pages?slug=current-news", fullPath: "/api/v1/related-pages"},
		{url: "/api/v1/pages/current-news", fullPath: "/api/v1/pages/*slug"},
	}
	for _, tt := range tests {
		t.Run(tt.fullPath, func(t *testing.T) {
			matchedPath = ""
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNoContent {
				t.Fatalf("expected matched route status 204, got %d", w.Code)
			}
			if matchedPath != tt.fullPath {
				t.Fatalf("expected route %q, got %q", tt.fullPath, matchedPath)
			}
		})
	}
}
