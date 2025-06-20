package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/auth/dto/request"
	"github.com/gogaruda/auth/auth/service"
	"github.com/gogaruda/auth/pkg/apperror"
	"github.com/gogaruda/auth/pkg/response"
	"github.com/gogaruda/auth/pkg/validates"
	"strconv"
)

type UserHandler struct {
	service   service.UserService
	Validator *validates.Validates
}

func NewUserHandler(s service.UserService, v *validates.Validates) *UserHandler {
	return &UserHandler{service: s, Validator: v}
}

// GetAllUsers godoc
// @Summary Get all users with pagination
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} apperror.InitError
// @Router /api/auth/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := h.service.GetAll(limit, offset)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	meta := response.MetaData{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	response.OK(c, users, "query ok", &meta)
}

// CreateUser godoc
// @Summary Create a new user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.CreateUserRequest true "User data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} apperror.InitError
// @Failure 401 {object} apperror.InitError
// @Router /api/auth/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if !h.Validator.ValidateJSON(c, &req) {
		return
	}

	if err := h.service.Create(req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.Created(c, nil, "query ok")
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} apperror.InitError
// @Failure 404 {object} apperror.InitError
// @Router /api/auth/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.service.GetByID(userID)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, user, "query ok", nil)
}

// UpdateUser godoc
// @Summary Update user by ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param data body request.UpdateUserRequest true "Updated user data"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} apperror.InitError
// @Failure 401 {object} apperror.InitError
// @Failure 404 {object} apperror.InitError
// @Router /api/auth/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userID = c.Param("id")
	var req request.UpdateUserRequest
	if !h.Validator.ValidateJSON(c, &req) {
		return
	}

	if err := h.service.UpdateUser(userID, req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.Created(c, nil, "query ok")
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 204 {object} nil
// @Failure 401 {object} apperror.InitError
// @Failure 404 {object} apperror.InitError
// @Router /api/auth/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.Delete(userID); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.NoContent(c)
}
