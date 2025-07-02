package middleware

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type middleware struct {
	db  *sql.DB
	cfg config.JWTConfig
}

func NewMiddleware(d *sql.DB, c config.JWTConfig) Middleware {
	return &middleware{db: d, cfg: c}
}

func (m *middleware) AuthMiddleware() gin.HandlerFunc {
	secret := m.cfg.Secret
	if secret == "" {
		panic("JWT_SECRET tidak ditemukan di file .env")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "header authorized tidak ditemukan")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "format authorization harus: bearer {token}")
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "token tidak valid atau sudah kadaluarsa")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "klaim token tidak valid")
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			response.Unauthorized(c, "token sudah kadaluarsa")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			response.Unauthorized(c, "token tidak memiliki identitas pengguna (user_id)")
			return
		}

		tokenVersion, ok := claims["token_version"].(string)
		if !ok {
			response.Unauthorized(c, "token version tidak valid")
			return
		}

		var user model.UserModel
		if err := m.db.QueryRow(`SELECT token_version FROM users WHERE id = ?`, userID).Scan(&user.TokenVersion); err != nil {
			response.Unauthorized(c, "user tidak ditemukan")
			return
		}

		if user.TokenVersion != tokenVersion {
			response.Unauthorized(c, "token sudah tidak berlaku, silahkan login lagi!")
			return
		}

		rolesInterface, ok := claims["roles"].([]interface{})
		if !ok {
			response.Unauthorized(c, "format roles dalam token tidak valid")
			return
		}

		roles := make([]string, 0, len(rolesInterface))
		for _, r := range rolesInterface {
			roleStr, ok := r.(string)
			if !ok {
				response.Unauthorized(c, "role tidak valid")
				return
			}
			roles = append(roles, roleStr)
		}

		c.Set("user_id", userID)
		c.Set("roles", roles)

		c.Next()
	}
}
