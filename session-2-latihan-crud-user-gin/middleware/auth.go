package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(401, gin.H{"error": "Authorization basic token required"})
			c.Abort()
			return
		}

		const (
			expectedusername = "admin"
			expectedpassword = "admin1234"
		)

		isValid := username == expectedusername && password == expectedpassword
		if !isValid {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		c.Next()

	}
}
