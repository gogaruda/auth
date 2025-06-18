package container

import (
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/system/config"
)

type AppService struct {
	UserService service.UserService
	AuthService service.AuthService
}

func InitApp() *AppService {
	config.LoadENV()
	config.ConnectDB()
	db := config.DB

	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db)

	authService := service.NewAuthService(authRepository)
	userService := service.NewUserService(userRepository)

	return &AppService{
		UserService: userService,
		AuthService: authService,
	}
}
