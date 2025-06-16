package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJWT(userID, tokenVersion string, roles []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":       userID,
		"token_version": tokenVersion,
		"roles":         roles,
		"exp":           now.Add(time.Hour * 24).Unix(),
		"iat":           now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}
