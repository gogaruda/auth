package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sql/internal/handler"
	"sql/internal/middleware"
	"sql/pkg/system/container"
)

func InitRouter(r *gin.Engine, app *container.AppService) {
	r.Use(middleware.CORSMiddleware())

	authHandler := handler.NewAuthHandler(app.AuthService)

	api := r.Group("/api")
	api.POST("/login", authHandler.Login)

	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/coba-auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Yee, Anda berhasil login!"})
	})
}
