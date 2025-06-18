package apperror

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type HTTPErrorMap struct {
	Code        string
	Status      int
	UserMessage string
}

var defaultErrorMap = []HTTPErrorMap{
	{CodeUserNotFound, http.StatusNotFound, "User tidak ditemukan"},
	{CodeResourceNotFound, http.StatusNotFound, "Data tidak ditemukan"},
	{CodeInvalidInput, http.StatusBadRequest, "Input tidak valid"},
	{CodeUserConflict, http.StatusConflict, "User sudah ada"},
	{CodeResourceConflict, http.StatusConflict, "Data bentrok atau duplikat"},
	{CodeUnauthorized, http.StatusUnauthorized, "Anda harus login"},
	{CodeForbidden, http.StatusForbidden, "Akses ditolak"},
	{CodeNotImplemented, http.StatusNotImplemented, "Fitur belum tersedia"},
	{CodeTimeout, http.StatusGatewayTimeout, "Permintaan melebihi waktu tunggu"},
	{CodeDependencyError, http.StatusBadGateway, "Kesalahan dari layanan eksternal"},
	{CodeDBNoRows, http.StatusNotFound, "Data tidak tersedia"},
	{CodeDBConstraint, http.StatusConflict, "Gagal menyimpan data: constraint"},
	{CodeDBTxFailed, http.StatusInternalServerError, "Transaksi database gagal"},
	{CodeDBConnFailed, http.StatusServiceUnavailable, "Koneksi database gagal"},
	{CodeDBError, http.StatusInternalServerError, "Kesalahan database"},
	{CodeInternalError, http.StatusInternalServerError, "Terjadi kesalahan internal"},
	{CodeRoleNotFound, http.StatusBadRequest, "Role tidak ditemukan"},
	{CodeAuthNotFound, http.StatusNotFound, "username/email atau password salah"},
}

func HandleHTTPError(c *gin.Context, err error) {
	var initErr *InitError
	if errors.As(err, &initErr) {
		// Cek log untuk semua error yang dikenali
		if initErr.Err != nil {
			log.Printf("[ERROR]: %s | DETAIL: %+v\n", initErr.Code, initErr.Err)
		} else {
			log.Printf("[ERROR]: %s | MESSAGE: %s\n", initErr.Code, initErr.Message)
		}

		for _, mapping := range defaultErrorMap {
			if mapping.Code == initErr.Code {
				c.JSON(mapping.Status, gin.H{
					"code":    mapping.Status,
					"status":  "error",
					"message": mapping.UserMessage,
				})
				return
			}
		}
	}

	// Fallback: tidak dikenali
	log.Printf("[UNHANDLED ERROR]: %v\n", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "terjadi kesalahan server",
	})
}
