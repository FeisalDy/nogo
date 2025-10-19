package middleware

import (
	"fmt"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware is a middleware for Casbin authorization
// It checks if the authenticated user has permission to access the resource
func CasbinMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := GetUserID(c)
		if !exists {
			utils.RespondWithAppError(c, errors.ErrAuthTokenMissing)
			c.Abort()
			return
		}

		// Get Casbin enforcer (singleton, already initialized)
		enforcer := casbinService.GetEnforcer()
		if enforcer == nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		// Check permission
		userSubject := fmt.Sprintf("user:%d", userID)
		allowed, err := enforcer.Enforce(userSubject, resource, action)
		if err != nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		if !allowed {
			utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

// DynamicCasbinMiddleware checks permissions based on route parameters
// The resource and action can be extracted from the request context
func DynamicCasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := GetUserID(c)
		if !exists {
			utils.RespondWithAppError(c, errors.ErrAuthTokenMissing)
			c.Abort()
			return
		}

		// Extract resource from path - e.g., /api/v1/users -> "users"
		// You can customize this based on your route structure
		resource := c.GetString("casbin_resource")
		action := c.GetString("casbin_action")

		if resource == "" || action == "" {
			// Fall back to method-based action
			switch c.Request.Method {
			case "GET":
				action = "read"
			case "POST":
				action = "write"
			case "PUT", "PATCH":
				action = "update"
			case "DELETE":
				action = "delete"
			default:
				action = "read"
			}

			// Extract resource from path
			if resource == "" {
				resource = c.FullPath()
			}
		}

		// Get Casbin enforcer
		enforcer := casbinService.GetEnforcer()
		if enforcer == nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		// Check permission
		userSubject := fmt.Sprintf("user:%d", userID)
		allowed, err := enforcer.Enforce(userSubject, resource, action)
		if err != nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		if !allowed {
			utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole checks if user has any of the specified roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			utils.RespondWithAppError(c, errors.ErrAuthTokenMissing)
			c.Abort()
			return
		}

		enforcer := casbinService.GetEnforcer()
		if enforcer == nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		userSubject := fmt.Sprintf("user:%d", userID)
		userRoles, err := enforcer.GetRolesForUser(userSubject)
		if err != nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range roles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllRoles checks if user has all of the specified roles
func RequireAllRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			utils.RespondWithAppError(c, errors.ErrAuthTokenMissing)
			c.Abort()
			return
		}

		enforcer := casbinService.GetEnforcer()
		if enforcer == nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		userSubject := fmt.Sprintf("user:%d", userID)
		userRoles, err := enforcer.GetRolesForUser(userSubject)
		if err != nil {
			utils.RespondWithAppError(c, errors.ErrInternalServer)
			c.Abort()
			return
		}

		// Convert user roles to map for faster lookup
		userRoleMap := make(map[string]bool)
		for _, role := range userRoles {
			userRoleMap[role] = true
		}

		// Check if user has all required roles
		for _, requiredRole := range roles {
			if !userRoleMap[requiredRole] {
				utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// SetCasbinResource sets the resource name for dynamic Casbin middleware
func SetCasbinResource(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("casbin_resource", resource)
		c.Next()
	}
}

// SetCasbinAction sets the action name for dynamic Casbin middleware
func SetCasbinAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("casbin_action", action)
		c.Next()
	}
}

// SetCasbinParams is a helper to set both resource and action
func SetCasbinParams(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("casbin_resource", resource)
		c.Set("casbin_action", action)
		c.Next()
	}
}

// Helper to get common action from HTTP method
func GetActionFromMethod(method string) string {
	switch method {
	case "GET", "HEAD":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "read"
	}
}

// PermissionChecker is a helper function to check permissions without middleware
func PermissionChecker(c *gin.Context, resource, action string) (bool, error) {
	userID, exists := GetUserID(c)
	if !exists {
		return false, fmt.Errorf("user not authenticated")
	}

	enforcer := casbinService.GetEnforcer()
	if enforcer == nil {
		return false, fmt.Errorf("casbin not initialized")
	}

	userSubject := fmt.Sprintf("user:%d", userID)
	return enforcer.Enforce(userSubject, resource, action)
}
