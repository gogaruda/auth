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
	GoogleAuthService        service.GoogleAuthService
}

func InitBootstrap(db *sql.DB, config *config.AppConfig) *Service {
	mail := mailer.NewMailer(config.Mail)
	ut := utils.NewUtils(config)

	userRepo := repository.NewUserRepository(db)
	emailRepo := repository.NewEmailVerificationRepository(db)
	authRepo := repository.NewAuthRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	userService := service.NewUserService(userRepo)
	emailService := service.NewEmailVerificationService(emailRepo, mail, ut, config.Mail, userService)
	authService := service.NewAuthService(authRepo, roleRepo, config, ut, emailService)
	googleService := service.NewGoogleAuthService(userRepo, roleRepo, authRepo, config, ut)

	newMiddleware := middleware.NewMiddleware(db, config.JWT, config.Cors)
	return &Service{
		AuthService:              authService,
		Middleware:               newMiddleware,
		EmailVerificationService: emailService,
		GoogleAuthService:        googleService,
	}
}
