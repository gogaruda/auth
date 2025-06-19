package auth

import (
	"database/sql"
	"github.com/gogaruda/auth/auth/repository"
	"github.com/gogaruda/auth/auth/service"
)

type Module struct {
	UserService service.UserService
	AuthService service.AuthService
}

func InitAuthModule(db *sql.DB) *Module {
	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db)

	authService := service.NewAuthService(authRepository)
	userService := service.NewUserService(userRepository)

	return &Module{
		UserService: userService,
		AuthService: authService,
	}
}
