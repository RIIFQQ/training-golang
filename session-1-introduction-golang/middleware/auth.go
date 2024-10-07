package middleware

import "github.com/gin-gonic/gin"

//membuat function AuthMiddleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization") //ambil header Authorization

		if token == "" { //jika token kosong
			c.JSON(401, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		if token != "token-rahasia" { //jika token tidak sama dengan "token
			c.JSON(401, gin.H{"error": "Invalid Authorization token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
