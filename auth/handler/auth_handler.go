package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/auth/dto/request"
	"github.com/gogaruda/auth/auth/service"
	"github.com/gogaruda/auth/pkg/apperror"
	"github.com/gogaruda/auth/pkg/response"
	"github.com/gogaruda/auth/pkg/validates"
	"net/http"
)

type AuthHandler struct {
	service   service.AuthService
	Validator *validates.Validates
}

func NewAuthHandler(h service.AuthService, v *validates.Validates) *AuthHandler {
	return &AuthHandler{service: h, Validator: v}
}

// Login godoc
// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.AuthLoginRequest true "Login data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} apperror.InitError
// @Failure 401 {object} apperror.InitError
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.AuthLoginRequest
	if !h.Validator.ValidateJSON(c, &req) {
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, nil, token, nil)
}

// Register godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.AuthRegisterRequest true "Register data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} apperror.InitError
// @Failure 409 {object} apperror.InitError
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.AuthRegisterRequest
	if !h.Validator.ValidateJSON(c, &req) {
		return
	}

	if err := h.service.Register(req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.Created(c, nil, "query ok")
}

// Logout godoc
// @Summary Logout current user
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} apperror.InitError
// @Router /api/auth/logout [post]
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
