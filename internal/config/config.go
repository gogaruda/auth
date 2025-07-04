package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
	Mode   GinModeConfig
	Cors   CORSConfig
	Mail   EmailConfig
	Google *oauth2.Config
}

func LoadConfig() *AppConfig {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return &AppConfig{
		DB: DBConfig{
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:         os.Getenv("JWT_SECRET"),
			AccessTokenTTL: 15 * time.Minute,
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
		Mode: GinModeConfig{
			Debug: getModeOrDefault("GIN_MODE", "debug"),
		},
		Cors: CORSConfig{
			AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
			AllowMethods:     strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ","),
			AllowHeaders:     strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ","),
			AllowCredentials: os.Getenv("CORS_ALLOW_CREDENTIALS") == "true",
		},
		Mail: EmailConfig{
			MailHost:        os.Getenv("MAIL_HOST"),
			MailPort:        port,
			MailUsername:    os.Getenv("MAIL_USERNAME"),
			MailPassword:    os.Getenv("MAIL_PASSWORD"),
			MailFromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
			FrontVerifyUrl:  os.Getenv("FRONTEND_VERIFY_URL"),
		},
		Google: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}
