package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"training-golang/session-3-unit-test/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHelloMessage(t *testing.T) {
	t.Run("Positive Case - Correct Message", func(t *testing.T) {
		expectedOutput := "Hallo dari Gin!"
		actualOutput := handler.GetHelloMessage()
		require.Equal(t, expectedOutput, actualOutput, "The message should be '%s'", expectedOutput)
	})
}

func TestRootHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/", handler.RootHandler)

	//create a new http request
	req, _ := http.NewRequest("GET", "/", nil)

	//create a responserecorder to record the response
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"Hallo dari Gin!"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

type JsonRequest struct {
	Message string `json:"message"`
}

func TestPostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/", handler.PostHandler)

	t.Run("Positive Case - Correct Message", func(t *testing.T) {
		//persiapan data JSON
		requestBody := JsonRequest{Message: "Hello from test!"}
		requestBodyBytes, _ := json.Marshal(requestBody)

		//buat permintaan HTTP Post
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		//buat response recorder untuk merekam respon
		w := httptest.NewRecorder()

		//lakukan permintaan
		//fmt.println()
		//penamaan route beardasarkan  router := gin.Default() dan router.POST("/", handler.PostHandler)
		router.ServeHTTP(w, req)

		//cek kode status
		assert.Equal(t, http.StatusOK, w.Code)

		//cek body response
		expectedBody := `{"message":"Hello from test!"}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})

	t.Run("Negative Case - EOF Error", func(t *testing.T) {
		// persiapan data Json yang salah
		requestBody := ""
		requestBodyBytes := []byte(requestBody)

		// buat permintaan HTTP Post
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		//cek kode status
		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Contains(t, w.Body.String(), "{\"error\":\"EOF\"}")
	})
}
