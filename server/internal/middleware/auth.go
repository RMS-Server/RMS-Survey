package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rms-survey/server/internal/config"
	jwtpkg "github.com/rms-survey/server/internal/pkg/jwt"
	"github.com/rms-survey/server/internal/pkg/response"
)

// publicPrefixes lists API path prefixes that skip authentication.
var publicPrefixes = []string{
	"/api/public/",
	"/api/oauth/authorize",
	"/api/oauth/callback",
	"/captcha/",
}

// isPublicPath returns true if the request path should bypass auth.
func isPublicPath(method, path string) bool {
	// Only /api/* paths need auth check
	if !strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/captcha/") {
		return true
	}

	for _, p := range publicPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	// GET /api/system is public
	if method == "GET" && path == "/api/system" {
		return true
	}
	// GET /api/file/** is public
	if method == "GET" && strings.HasPrefix(path, "/api/file/") {
		return true
	}
	return false
}

// Auth extracts and validates JWT from cookie, Authorization header, or query param.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicPath(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		var tokenStr string

		// Try cookie first
		if cookie, err := c.Cookie(config.C.JWT.CookieName); err == nil && cookie != "" {
			tokenStr = cookie
		}

		// Try Authorization header (Bearer token)
		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// Support both "Bearer <token>" and just "<token>"
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
				} else {
					tokenStr = authHeader
				}
			}
		}

		// Fall back to query param
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		user, err := jwtpkg.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
