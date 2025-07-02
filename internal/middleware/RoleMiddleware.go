package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/pkg/response"
	"strings"
)

type RoleMatchType int

const (
	MatchAll RoleMatchType = iota
	MatchAny
)

func (m *middleware) RoleMiddleware(matchType RoleMatchType, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("roles")
		if !ok {
			response.Unauthorized(c, "role tidak ditemukan pada context")
			return
		}

		userRoles, ok := val.([]string)
		if !ok {
			response.ServerError(c, "format roles tidak valid")
			return
		}

		if !matchRoles(userRoles, requiredRoles, matchType) {
			response.Forbidden(c, "Anda tidak berhak mengakses halaman ini")
			return
		}

		c.Next()
	}
}

func matchRoles(userRoles, requiredRoles []string, matchType RoleMatchType) bool {
	roleSet := make(map[string]struct{}, len(userRoles))
	for _, role := range userRoles {
		roleSet[strings.ToLower(role)] = struct{}{}
	}

	matchCount := 0
	for _, required := range requiredRoles {
		if _, ok := roleSet[strings.ToLower(required)]; ok {
			matchCount++
			if matchType == MatchAny {
				return true
			}
		} else if matchType == MatchAll {
			return false
		}
	}

	return matchType == MatchAll && matchCount == len(requiredRoles)
}
