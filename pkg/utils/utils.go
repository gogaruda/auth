package utils

import "github.com/gogaruda/auth/internal/config"

type Utils interface {
	GenerateJWT(userID, tokenVersion string, isVerified bool, roles []string, cfg *config.AppConfig) (string, error)
	GenerateULID() string
	GenerateHash(password string) (string, error)
	CompareHash(hash, password string) bool
	GenerateUsernameFromName(name string) string
}

type utils struct {
	config *config.AppConfig
}

func NewUtils(cfg *config.AppConfig) Utils {
	return &utils{config: cfg}
}
