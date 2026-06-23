package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionResolver interface {
	EffectivePermissions(userID uint64) ([]string, error)
}

func LoadPermissions(resolver PermissionResolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint64)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
			c.Abort()
			return
		}

		permissions, err := resolver.EffectivePermissions(userID)
		if err != nil {
			log.Printf("failed to load permissions for user %d: %v", userID, err)
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
			c.Abort()
			return
		}

		c.Set("permissions", permissions)
		c.Next()
	}
}

func RBAC(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hasAnyRequiredPermission(c, []string{requiredPermission}) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
		c.Abort()
	}
}

func RBACAny(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hasAnyRequiredPermission(c, requiredPermissions) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "forbidden"})
		c.Abort()
	}
}

func hasAnyRequiredPermission(c *gin.Context, requiredPermissions []string) bool {
	if len(requiredPermissions) == 0 {
		return true
	}

	permissionsValue, exists := c.Get("permissions")
	if !exists {
		return false
	}

	permissions, ok := permissionsValue.([]string)
	if !ok {
		return false
	}

	for _, permission := range permissions {
		for _, requiredPermission := range requiredPermissions {
			if permission == requiredPermission {
				return true
			}
		}
	}
	return false
}
