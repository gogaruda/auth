package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/auth"
	"github.com/gogaruda/auth/auth/config"
	"github.com/gogaruda/auth/auth/middleware"
	_ "github.com/gogaruda/auth/docs"
	"github.com/gogaruda/auth/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// Swagger documentation
// @title Auth - REST API Docs
// @description Auth system
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadENV()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := config.ConnectDB()
	app := auth.InitAuthModule(db)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	swagger.RegisterSwaggerRoutes(api.Group("/auth"))
	auth.RegisterAuthRoutes(api.Group("/auth"), app.AuthService, app.UserService)

	port := os.Getenv("APP_PORT")
	fmt.Println(port)
	if port == "" {
		port = "8080"
	}

	_ = r.Run(":" + port)
}
