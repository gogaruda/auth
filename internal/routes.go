package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/bootstrap"
	"github.com/gogaruda/auth/internal/handler"
	"github.com/gogaruda/valigo"
	"net/http"
)

func RouteRegister(r *gin.Engine, app *bootstrap.Service) {
	v := valigo.NewValigo()

	authHandler := handler.NewAuthHandler(app.AuthService, v)

	api := r.Group("/api")
	api.POST("/login", authHandler.Login)

	auth := api.Group("/")
	auth.Use(app.Middleware.AuthMiddleware())
	auth.GET("/coba", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  "success",
			"meesage": "selamat datang di auth middleware",
		})
	})
	auth.POST("/logout", authHandler.Logout)
}
