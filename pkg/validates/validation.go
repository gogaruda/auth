package validates

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gogaruda/auth/pkg/response"
)

type Validates struct {
	validator *validator.Validate
}

func NewValidates(v *validator.Validate) *Validates {
	return &Validates{validator: v}
}

// Ambil old data
type Sanitizable interface {
	Sanitize() map[string]any
}

func (vu *Validates) ValidateJSON(c *gin.Context, input any) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorsMap := buildValidationErrorMap(input, verr)

			var old any
			if s, ok := input.(Sanitizable); ok {
				old = s.Sanitize()
			} else {
				old = input
			}

			response.BadRequest(c, gin.H{
				"errors": errorsMap,
				"old":    old,
			}, "Validasi gagal")
			return false
		}
		response.BadRequest(c, gin.H{
			"error": "Format data tidak valid. Pastikan format JSON benar.",
		}, "Permintaan tidak valid")
		return false
	}
	return true
}

// Validasi username/email dan lain sebagainya
func (vu *Validates) ValidateBussiness(c *gin.Context, input Sanitizable, FieldErrors map[string]string) bool {
	if len(FieldErrors) == 0 {
		return true
	}

	response.BadRequest(c, gin.H{
		"errors": FieldErrors,
		"old":    input.Sanitize(),
	}, "validasi gagal")

	return false
}
