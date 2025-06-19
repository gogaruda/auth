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

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.service.GetByID(userID)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, user, "query ok", nil)
}

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

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.Delete(userID); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.NoContent(c)
}
