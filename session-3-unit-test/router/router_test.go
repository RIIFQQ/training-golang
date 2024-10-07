package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"training-golang/session-3-unit-test/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter_RootHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	//buat permintaan HTTP POST dengan JSON yang tidak valid

	req, _ := http.NewRequest("GET", "/", nil)

	//jalankan permintaan HTTP
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//periksa kode status
	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"Hallo dari Gin!"}`
	assert.JSONEq(t, expectedBody, w.Body.String())

}

func TestPostHandler_PositiveCase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	//buat permintaan HTTP POST dengan JSON yang tidak valid
	requestBody := map[string]string{"message": "Test Message"}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/v1/api/post", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Authorization", "token-rahasia")
	req.Header.Set("Content-Type", "application/json")

	//jalankan permintaan HTTP
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//periksa kode status
	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"Test Message"}`
	assert.JSONEq(t, expectedBody, w.Body.String())

}

func TestPostHandler_NegativeCase_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	//buat permintaan HTTP POST dengan JSON yang tidak valid
	req, _ := http.NewRequest("POST", "/v1/api/post", bytes.NewBufferString("{Invalid JSON}"))
	req.Header.Set("Authorization", "token-rahasia")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	//jalankan permintaan HTTP
	r.ServeHTTP(w, req)

	//periksa kode status
	assert.Equal(t, http.StatusBadRequest, w.Code)

	//periksa body response
	assert.Contains(t, w.Body.String(), "invalid character")

}

func TestPostHandler_NegativeCase_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	//buat permintaan HTTP POST dengan JSON yang tidak valid
	req, _ := http.NewRequest("POST", "/v1/api/post", nil)

	w := httptest.NewRecorder()

	//jalankan permintaan HTTP
	r.ServeHTTP(w, req)

	//periksa kode status
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	//periksa body response
	assert.Contains(t, w.Body.String(), "{\"error\":\"Authorization token required\"}")

}
