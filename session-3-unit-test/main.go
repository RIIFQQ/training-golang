package main

import (
	"training-golang/session-3-unit-test/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	router.SetupRouter(r)

	r.Run(":8080")
}
