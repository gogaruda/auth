package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gogaruda/auth/internal/handler"
	"github.com/gogaruda/auth/internal/middleware"
	"github.com/gogaruda/auth/pkg/validates"
	"github.com/gogaruda/auth/system/container"
	"net/http"
)

func InitRouter(r *gin.Engine, app *container.AppService) {
	r.Use(middleware.CORSMiddleware())

	v := validator.New()
	valid := validates.NewValidates(v)

	authHandler := handler.NewAuthHandler(app.AuthService, valid)
	userHandler := handler.NewUserHandler(app.UserService, valid)

	api := r.Group("/api")
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)

	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/coba-auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Yee, Anda berhasil login!"})
	})

	superAdmin := auth.Group("/")
	superAdmin.Use(middleware.RoleMiddleware(middleware.MatchAny, "super-admin"))
	{
		superAdmin.GET("/users", userHandler.GetAllUsers)
		superAdmin.POST("/users", userHandler.CreateUser)
		superAdmin.GET("/users/:id", userHandler.GetUserByID)
		superAdmin.PUT("/users/:id", userHandler.UpdateUser)
		superAdmin.DELETE("/users/:id", userHandler.DeleteUser)
	}
}
