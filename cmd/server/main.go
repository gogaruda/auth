package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/auth"
	"github.com/gogaruda/auth/auth/config"
	"github.com/gogaruda/pkg/middleware"
	"os"
	"strings"
)

func getAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:3000"}
	}
	return strings.Split(origins, ",")
}

func main() {
	config.LoadENV()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := config.ConnectDB()
	app := auth.InitAuthModule(db)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(getAllowedOrigins()))

	api := r.Group("/api")
	auth.RegisterAuthRoutes(api.Group("/auth"), app.AuthService, app.UserService)

	port := os.Getenv("APP_PORT")
	fmt.Println(port)
	if port == "" {
		port = "8080"
	}

	_ = r.Run(":" + port)
}
