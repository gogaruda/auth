package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth"
	"net/http"
)

func RegisterSwaggerRoutes(r *gin.RouterGroup) {
	r.GET("/swagger.yaml", func(c *gin.Context) {
		data, err := auth.SwaggerFS.ReadFile("docs/swagger.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "swagger file not found"})
			return
		}
		c.Data(http.StatusOK, "application/x-yaml", data)
	})
}
