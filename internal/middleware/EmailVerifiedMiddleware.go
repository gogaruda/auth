package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/pkg/response"
)

func (m *middleware) EmailVerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.NewResponder(c)
		isVerifiedRaw, exists := c.Get("is_verified")
		if !exists {
			res.Forbidden("status verifikasi email tidak tersedia")
			return
		}

		isVerified, ok := isVerifiedRaw.(bool)
		if !ok || !isVerified {
			res.Forbidden("email belum diverifikasi")
			return
		}

		c.Next()
	}
}
