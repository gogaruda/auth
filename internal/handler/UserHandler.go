package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/service"
	"github.com/gogaruda/auth/pkg/response"
	"github.com/gogaruda/valigo"
	"strconv"
)

type UserHandler struct {
	service service.UserService
	valid   *valigo.Valigo
}

func NewUserHandler(s service.UserService, v *valigo.Valigo) *UserHandler {
	return &UserHandler{service: s, valid: v}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	res := response.NewResponder(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := h.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	meta := response.MetaData{
		Total: total,
		Page:  page,
		Limit: limit,
	}

	res.OK(users, "query ok", &meta)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.UserCreateRequest
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.Created(nil, "user berhasil dibuat")
}
