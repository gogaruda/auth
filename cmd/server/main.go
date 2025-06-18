package main

import (
	"github.com/gin-gonic/gin"
	"sql/system/container"
	"sql/system/router"
)

func main() {
	r := gin.Default()

	app := container.InitApp()
	router.InitRouter(r, app)
	r.Run()
}
