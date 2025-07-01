package bootstrap

import (
	"database/sql"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/utils"
)

type Service struct {
	AuthService service.AuthService
}

func InitBootstrap(db *sql.DB, config *config.AppConfig) *Service {
	hasher := utils.NewBcryptHasher()
	jwt := utils.NewJWTGenerated()
	id := utils.NewULIDCreate()

	authRepo := repository.NewAuthRepository(db, id)

	authService := service.NewAuthService(authRepo, config, hasher, jwt)

	return &Service{
		AuthService: authService,
	}
}
