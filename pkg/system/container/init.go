package container

import (
	"sql/internal/repository"
	"sql/internal/service"
	"sql/pkg/system/config"
)

type AppService struct {
	UserService service.UserService
}

func InitApp() *AppService {
	config.LoadENV()
	config.ConnectDB()
	db := config.DB

	userRepository := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepository)

	return &AppService{
		UserService: userService,
	}
}
