package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/response"
	"github.com/gogaruda/valigo"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
	valid       *valigo.Valigo
}

func NewAuthHandler(a service.AuthService, v *valigo.Valigo) *AuthHandler {
	return &AuthHandler{authService: a, valid: v}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"status": "success",
		"token":  token,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.authService.Logout(userID.(string)); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, nil, "logout success!", nil)
}
