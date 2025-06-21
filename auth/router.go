package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gogaruda/auth/auth/handler"
	"github.com/gogaruda/auth/auth/middleware"
	"github.com/gogaruda/auth/auth/service"
	"github.com/gogaruda/pkg/validates"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authService service.AuthService, userService service.UserService) {
	v := validator.New()
	valid := validates.NewValidates(v)

	authHandler := handler.NewAuthHandler(authService, valid)
	userHandler := handler.NewUserHandler(userService, valid)

	rg.POST("/login", authHandler.Login)
	rg.POST("/register", authHandler.Register)

	auth := rg.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/logout", authHandler.Logout)

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
}
