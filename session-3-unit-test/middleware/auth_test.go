package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"training-golang/session-3-unit-test/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_PositifCase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	//inisialisasi router
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	//handler yang hanya dapat diakses dengan token
	r.GET("/private", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Private Data"})
	})

	//buat permintaan HTTP Get ke endpoint /private dengan token yang valid
	req, _ := http.NewRequest("GET", "/private", nil)
	req.Header.Set("Authorization", "token-rahasia")

	w := httptest.NewRecorder()

	//lakukan permintaan HTTP
	r.ServeHTTP(w, req)

	//periksa status code
	assert.Equal(t, http.StatusOK, w.Code)

	//periksa body response
	assert.JSONEq(t, `{"message":"Private Data"}`, w.Body.String())

}

func TestAuthMidlleware_Negative_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	//inisialisasi router
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	//handler yang hanya dapat diakses dengan token
	r.GET("/private", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Private Data"})
	})

	//buat permintaan HTTP Get ke endpoint /private dengan token yang valid
	req, _ := http.NewRequest("GET", "/private", nil)

	w := httptest.NewRecorder()

	//lakukan permintaan HTTP
	r.ServeHTTP(w, req)

	//periksa status code
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	//periksa body response
	//assert.Equal(t, `{"error":"Authorization token required"}`, w.Body.String())
	assert.Contains(t, w.Body.String(), "Authorization token required")
}

func TestAuthMiddleware_PositifCase_(t *testing.T) {
	gin.SetMode(gin.TestMode)
	//inisialisasi router
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	//handler yang hanya dapat diakses dengan token
	r.GET("/private", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Private Data"})
	})

	//buat permintaan HTTP Get ke endpoint /private dengan token yang valid
	req, _ := http.NewRequest("GET", "/private", nil)
	req.Header.Set("Authorization", "token-terbuka")

	w := httptest.NewRecorder()

	//lakukan permintaan HTTP
	r.ServeHTTP(w, req)

	//periksa status code
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	//periksa body response
	assert.Contains(t, w.Body.String(), "Invalid Authorization token")
}
