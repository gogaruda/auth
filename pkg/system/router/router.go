package router

import (
	"github.com/gin-gonic/gin"
	"sql/internal/handler"
	"sql/pkg/system/container"
)

func InitRouter(r *gin.Engine, app *container.AppService) {
	authHandler := handler.NewAuthHandler(app.AuthService)
	userHandler := handler.NewUserHandler(app.UserService)
	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/user", userHandler.GetAllUsers)
		api.GET("/user/:id", userHandler.GetByID)
	}
}
