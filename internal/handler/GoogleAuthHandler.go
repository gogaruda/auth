package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/service"
	"net/http"
)

type GoogleAuthHandler struct {
	service service.GoogleAuthService
}

func NewGoogleAuthHandler(s service.GoogleAuthService) *GoogleAuthHandler {
	return &GoogleAuthHandler{service: s}
}

func (h *GoogleAuthHandler) GoogleLogin(c *gin.Context) {
	url := h.service.Login()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *GoogleAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := h.service.Callback(c.Request.Context(), code)
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
