package middleware

import (
	"majoo-case1-rest-api/config"
	httpx "majoo-case1-rest-api/internal/http"
	"majoo-case1-rest-api/internal/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cookie-only auth using JWT
func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	secret := []byte(cfg.JWTSecret)
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			httpx.RespondWithError(c, http.StatusUnauthorized, "Authentication cookie required")
			c.Abort()
			return
		}
		claims, err := security.ValidateToken(secret, token)
		if err != nil {
			httpx.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Next()
	}
}
