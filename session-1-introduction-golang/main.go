package main

import (
	"training-golang/session-1-introduction-golang/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	router.SetupRouter(r) //memanggil fungsi SetupRouter dari package router

	r.Run(":8080")
}
