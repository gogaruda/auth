package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/apperror"
	"github.com/gogaruda/auth/pkg/response"
	"net/http"
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
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, nil, token, nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.AuthRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.Created(c, nil, "query ok")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id tidak ditemukan!"})
		return
	}

	if err := h.service.Logout(userID.(string)); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.NoContent(c)
}
