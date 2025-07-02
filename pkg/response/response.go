package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OK(c *gin.Context, data interface{}, message string, meta *MetaData) {
	response := APIResponse{
		Code:    http.StatusOK,
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	}

	if meta != nil {
		response.Meta = *meta
	}

	c.JSON(http.StatusOK, response)
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Code:    http.StatusCreated,
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func BadRequest(c *gin.Context, errors interface{}, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
		Code:    http.StatusBadRequest,
		Status:  StatusError,
		Message: message,
		Errors:  errors,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
		Code:    http.StatusUnauthorized,
		Status:  StatusError,
		Message: message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
		Code:    http.StatusForbidden,
		Status:  StatusError,
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, APIResponse{
		Code:    http.StatusNotFound,
		Status:  StatusError,
		Message: message,
	})
}

func ServerError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
		Code:    http.StatusInternalServerError,
		Status:  StatusError,
		Message: message,
	})
}
