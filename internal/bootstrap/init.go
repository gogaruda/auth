package bootstrap

import (
	"database/sql"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/middleware"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/utils"
)

type Service struct {
	AuthService service.AuthService
	Middleware  middleware.Middleware
}

func InitBootstrap(db *sql.DB, config *config.AppConfig) *Service {
	hasher := utils.NewBcryptHasher()
	jwt := utils.NewJWTGenerated()
	id := utils.NewULIDCreate()

	authRepo := repository.NewAuthRepository(db, id)

	authService := service.NewAuthService(authRepo, config, hasher, jwt, id)

	newMiddleware := middleware.NewMiddleware(db, config.JWT, config.Cors)
	return &Service{
		AuthService: authService,
		Middleware:  newMiddleware,
	}
}
