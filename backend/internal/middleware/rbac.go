package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var rolePermissions = map[string][]string{
	"admin":  {"admin:read", "admin:write", "projects:write", "content:write", "leads:read"},
	"editor": {"admin:read", "content:write", "leads:read"},
	"viewer": {"admin:read"},
}

func RBAC(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
			c.Abort()
			return
		}

		roleStr := role.(string)
		permissions, ok := rolePermissions[roleStr]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
			c.Abort()
			return
		}

		for _, p := range permissions {
			if p == requiredPermission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
		c.Abort()
	}
}
