package utils

import (
	"github.com/gogaruda/auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (u *utils) GenerateJWT(userID, tokenVersion string, isVerified bool, roles []string, cfg *config.AppConfig) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":       userID,
		"token_version": tokenVersion,
		"is_verified":   isVerified,
		"roles":         roles,
		"exp":           now.Add(cfg.JWT.AccessTokenTTL).Unix(),
		"iat":           now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := cfg.JWT.Secret

	return token.SignedString([]byte(secret))
}
