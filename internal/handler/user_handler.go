package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/apperror"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "Data berhasil diambil",
		"data":    users,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tambah data baru berhasil"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.service.GetByID(userID)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userID = c.Param("id")
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateUser(userID, req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Update berhasil!"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.Delete(userID); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
