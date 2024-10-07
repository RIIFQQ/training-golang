package middleware

import (
	"training-golang/session-4-unit-test-crud-user/config"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(401, gin.H{"error": "Authorization basic token required"})
			c.Abort()
			return
		}

		isValid := username == config.AuthBasicUsername && password == config.AuthBasicpassword
		if !isValid {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		c.Next()

	}
}
