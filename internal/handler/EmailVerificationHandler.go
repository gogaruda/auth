package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/response"
)

type EmailVerificationHandler struct {
	service service.EmailVerificationService
}

func NewEmailVerificationHandler(s service.EmailVerificationService) *EmailVerificationHandler {
	return &EmailVerificationHandler{service: s}
}

func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	res := response.NewResponder(c)
	token := c.Query("token")
	if token == "" {
		res.BadRequest(nil, "token tidak boleh kosong")
		return
	}

	if err := h.service.VerifyToken(c.Request.Context(), token); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(nil, "Email berhasil di verifikasi", nil)
}
