package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/bootstrap"
	"github.com/gogaruda/auth/internal/handler"
	"github.com/gogaruda/auth/internal/middleware"
	"github.com/gogaruda/valigo"
)

func RouteRegister(r *gin.Engine, app *bootstrap.Service) {
	v := valigo.NewValigo()

	authHandler := handler.NewAuthHandler(app.AuthService, v)
	emailHandler := handler.NewEmailVerificationHandler(app.EmailVerificationService)
	googleHandler := handler.NewGoogleAuthHandler(app.GoogleAuthService)
	userHandler := handler.NewUserHandler(app.UserService, v)

	r.Use(app.Middleware.CORSMiddleware())
	api := r.Group("/api")

	// auth
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	// email
	api.GET("/email-verification", emailHandler.VerifyEmail)

	// google OAuth2
	api.GET("/google/login", googleHandler.GoogleLogin)
	api.GET("/google/callback", googleHandler.GoogleCallback)

	auth := api.Group("")
	auth.Use(app.Middleware.AuthMiddleware())

	eVerify := auth.Group("")
	eVerify.Use(app.Middleware.EmailVerifiedMiddleware())

	// role super admin dan admin
	superAndAdmin := app.Middleware.RoleMiddleware(middleware.MatchAny, "super admin", "admin")

	// users
	eVerify.GET("/users", superAndAdmin, userHandler.GetAllUsers)
	eVerify.POST("/users/create", superAndAdmin, userHandler.CreateUser)

	// logout
	eVerify.POST("/logout", authHandler.Logout)
}
