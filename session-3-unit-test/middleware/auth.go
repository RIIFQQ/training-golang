package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// membuat function AuthMiddleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//periksa token disediakan atau tidak
		token := c.GetHeader("Authorization") //ambil header Authorization

		if token == "" { //jika token kosong
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		//verifikasi token
		if token != "token-rahasia" { //jika token tidak sama dengan "token
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization token"})
			c.Abort()
			return
		}
		//lanjutkan ke handler berikutnya jika token valid
		c.Next()
	}
}
