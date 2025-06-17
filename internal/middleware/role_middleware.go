package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type MatchMode string

const (
	MatchAny MatchMode = "any"
	MatchAll MatchMode = "all"
)

// Parameter:
// - matchMode: menentukan mode pencocokan antara peran pengguna dan peran yang dibutuhkan.
//   - MatchAny: akses diberikan jika pengguna memiliki setidaknya salah satu peran yang dibutuhkan.
//   - MatchAll: akses hanya diberikan jika pengguna memiliki semua peran yang dibutuhkan.
//
// - requiredRoles: daftar peran yang diwajibkan untuk mengakses endpoint.
//
// Contoh penggunaan:
//
//	r.GET("/admin", RoleMiddleware(MatchAny, "admin", "superadmin"), adminHandler)
//	r.POST("/manage", RoleMiddleware(MatchAll, "admin", "manager"), manageHandler)
//
// Jika pengguna tidak memenuhi syarat, middleware ini akan menghentikan eksekusi
// dan mengembalikan HTTP 403 dengan pesan kesalahan berbahasa Indonesia.
//
// Catatan:
// Pastikan authMiddleware() sudah menyimpan informasi roles ke dalam context sebelum middleware ini dijalankan.
//
// Contoh setup roles di context (misalnya dalam JWT middleware):
//
//	c.Set("roles", []string{"admin", "editor"})

func RoleMiddleware(matchMode MatchMode, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil roles dari context (harus di-set oleh middleware sebelumnya, seperti JWT)
		rolesInterface, exists := c.Get("roles")
		if !exists {
			log.Println("[RoleMiddleware] Peran pengguna tidak ditemukan dalam konteks")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Tidak memiliki izin akses"})
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			log.Println("[RoleMiddleware] Format peran pengguna tidak valid")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Tidak memiliki izin akses"})
			return
		}

		// Ubah ke map untuk pencarian cepat
		roleMap := make(map[string]bool)
		for _, role := range userRoles {
			roleMap[strings.ToLower(role)] = true
		}

		// Cocokkan berdasarkan mode
		switch matchMode {
		case MatchAll:
			for _, required := range requiredRoles {
				if !roleMap[strings.ToLower(required)] {
					log.Printf("[RoleMiddleware] Akses ditolak. Dibutuhkan SEMUA: %v | Dimiliki: %v\n", requiredRoles, userRoles)
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
					return
				}
			}
		case MatchAny:
			matched := false
			for _, required := range requiredRoles {
				if roleMap[strings.ToLower(required)] {
					matched = true
					break
				}
			}
			if !matched {
				log.Printf("[RoleMiddleware] Akses ditolak. Dibutuhkan SALAH SATU: %v | Dimiliki: %v\n", requiredRoles, userRoles)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
				return
			}
		default:
			log.Printf("[RoleMiddleware] Mode pencocokan tidak valid: %s\n", matchMode)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		c.Next()
	}
}
