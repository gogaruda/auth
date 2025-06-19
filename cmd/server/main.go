package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/auth"
	"github.com/gogaruda/auth/auth/config"
	"os"
)

func main() {
	config.LoadENV()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := config.ConnectDB()
	app := auth.InitAuthModule(db)

	r := gin.Default()
	api := r.Group("/api")
	auth.RegisterAuthRoutes(api.Group("/auth"), app.AuthService, app.UserService)

	port := os.Getenv("APP_PORT")
	fmt.Println(port)
	if port == "" {
		port = "8080"
	}

	_ = r.Run(":" + port)
}
