package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/dto/request"
	dto "github.com/gogaruda/auth/internal/dto/response"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/response"
	"github.com/gogaruda/valigo"
)

type AuthHandler struct {
	authService service.AuthService
	valid       *valigo.Valigo
}

func NewAuthHandler(a service.AuthService, v *valigo.Valigo) *AuthHandler {
	return &AuthHandler{authService: a, valid: v}
}

func (h *AuthHandler) Register(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.RegisterRequest
	req.Roles = []string{"tamu"}
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	if err := h.authService.Register(c.Request.Context(), req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.Created(nil, "registrasi berhasil")
}

func (h *AuthHandler) Login(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.LoginRequest
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(dto.LoginResponse{
		Token: token,
	}, "login berhasil", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	res := response.NewResponder(c)
	userID, _ := c.Get("user_id")
	if err := h.authService.Logout(userID.(string)); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(nil, "logout success!", nil)
}
