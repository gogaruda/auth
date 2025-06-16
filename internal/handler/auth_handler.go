package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sql/internal/dto/request"
	"sql/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(h service.AuthService) *AuthHandler {
	return &AuthHandler{service: h}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
