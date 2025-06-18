package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/system/container"
	"github.com/gogaruda/auth/system/router"
)

func main() {
	r := gin.Default()

	app := container.InitApp()
	router.InitRouter(r, app)
	r.Run()
}
