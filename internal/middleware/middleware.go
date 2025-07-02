package middleware

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/config"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	CORSMiddleware() gin.HandlerFunc
	RoleMiddleware(matchType RoleMatchType, requiredRoles ...string) gin.HandlerFunc
}

type middleware struct {
	db      *sql.DB
	cfg     config.JWTConfig
	corsCfg config.CORSConfig
}

func NewMiddleware(d *sql.DB, c config.JWTConfig, cc config.CORSConfig) Middleware {
	return &middleware{db: d, cfg: c, corsCfg: cc}
}
