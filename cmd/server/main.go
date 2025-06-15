package main

import (
	"github.com/gin-gonic/gin"
	"sql/pkg/system/container"
	"sql/pkg/system/router"
)

func main() {
	r := gin.Default()

	app := container.InitApp()
	router.InitRouter(r, app)
	r.Run()
}
