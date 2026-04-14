package middleware

import (
	"net/http"
	"strings"

	"taskflow/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		// Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token format",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// attach user info to context
		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("name", claims["name"])

		c.Next()
	}
}
