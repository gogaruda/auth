package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/system/config"
	"github.com/gogaruda/auth/system/container"
	"github.com/gogaruda/auth/system/router"
	"os"
)

func main() {
	config.LoadENV()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	app := container.InitApp()
	router.InitRouter(r, app)

	port := os.Getenv("APP_PORT")
	fmt.Println(port)
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
