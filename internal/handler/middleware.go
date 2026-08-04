package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"accounting/internal/auth"
)

// AuthMiddleware /login dışındaki tüm sayfa ve /api rotalarını oturum kontrolüne tabi tutar.
func AuthMiddleware(sessions *auth.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		token, err := c.Cookie(sessionCookie)
		if err != nil || token == "" {
			denyAuth(c)
			return
		}
		userID, ok := sessions.UserID(token)
		if !ok {
			denyAuth(c)
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}

func denyAuth(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
}
