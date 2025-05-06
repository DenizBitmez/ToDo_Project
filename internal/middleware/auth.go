package middleware

import (
	"ToDoProject/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		if claims.Role != requiredRole {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Forbidden"})
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
