package bootstrap

import (
	"database/sql"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/middleware"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/mailer"
	"github.com/gogaruda/auth/pkg/utils"
)

type Service struct {
	AuthService              service.AuthService
	Middleware               middleware.Middleware
	EmailVerificationService service.EmailVerificationService
}

func InitBootstrap(db *sql.DB, config *config.AppConfig) *Service {
	hasher := utils.NewBcryptHasher()
	jwt := utils.NewJWTGenerated()
	id := utils.NewULIDCreate()
	mail := mailer.NewMailer(config.Mail)

	userRepo := repository.NewUserRepository(db)
	emailRepo := repository.NewEmailVerificationRepository(db)
	authRepo := repository.NewAuthRepository(db)

	userService := service.NewUserService(userRepo)
	emailService := service.NewEmailVerificationService(emailRepo, mail, id, config.Mail, userService)
	authService := service.NewAuthService(authRepo, config, hasher, jwt, id, emailService)

	newMiddleware := middleware.NewMiddleware(db, config.JWT, config.Cors)
	return &Service{
		AuthService:              authService,
		Middleware:               newMiddleware,
		EmailVerificationService: emailService,
	}
}
