package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"sql/internal/model"
	"sql/pkg/system/config"
	"strings"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET tidak ditemukan di environment variable")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header Authorized tidak ditemukan"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format Authorization harus: bearer {token}"})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa!"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Klaim token tidak valid!"})
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token sudah kadaluarsa!"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak memiliki identitas pengguna (user_id)"})
			return
		}

		tokenVersion, ok := claims["token_version"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token version tidak valid!"})
			return
		}

		var user model.UserModel
		db := config.DB
		if err := db.QueryRow("SELECT token_version FROM users WHERE id = ?", userID).Scan(&user.TokenVersion); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan!"})
			return
		}

		fmt.Println("Hasil query: ", user.TokenVersion)
		fmt.Println("Dari hasil login: ", tokenVersion)

		if user.TokenVersion != tokenVersion {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token sudah tidak berlaku, silahkan login lagi!"})
			return
		}

		rolesInterface, ok := claims["roles"].([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format roles tidak valid dalam token"})
			return
		}

		roles := make([]string, 0, len(rolesInterface))
		for _, r := range rolesInterface {
			roleStr, ok := r.(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role tidak valid, harus string"})
				return
			}
			roles = append(roles, roleStr)
		}

		c.Set("user_id", userID)
		c.Set("roles", roles)
		c.Next()
	}
}
