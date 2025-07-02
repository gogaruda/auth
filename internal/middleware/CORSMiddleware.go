package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/config"
	"net/http"
	"strings"
)

func CORSMiddleware(corsCfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowed := false
		for _, o := range corsCfg.AllowOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsCfg.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsCfg.AllowHeaders, ","))

		if corsCfg.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
