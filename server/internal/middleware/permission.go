package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/dto"
	"github.com/rms-survey/server/internal/pkg/response"
)

// RequirePermission returns middleware that checks if the current user has the given permission.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("currentUser")
		if !exists {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		user, ok := val.(*dto.UserInfo)
		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		for _, role := range user.Roles {
			if role == permission || role == "admin" {
				c.Next()
				return
			}
		}

		response.Forbidden(c)
		c.Abort()
	}
}
