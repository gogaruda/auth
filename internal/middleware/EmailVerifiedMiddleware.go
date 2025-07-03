package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/pkg/response"
)

func (m *middleware) EmailVerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isVerifiedRaw, exists := c.Get("is_verified")
		if !exists {
			response.Forbidden(c, "status verifikasi email tidak tersedia")
			return
		}

		isVerified, ok := isVerifiedRaw.(bool)
		if !ok || !isVerified {
			response.Forbidden(c, "email belum diverifikasi")
			return
		}

		c.Next()
	}
}
