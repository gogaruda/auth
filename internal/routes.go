package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/bootstrap"
	"github.com/gogaruda/auth/internal/handler"
	"github.com/gogaruda/valigo"
)

func RouteRegister(r *gin.Engine, app *bootstrap.Service) {
	v := valigo.NewValigo()

	auth := handler.NewAuthHandler(app.AuthService, v)

	api := r.Group("/api")
	api.POST("/login", auth.Login)
}
