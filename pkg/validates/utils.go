package validates

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func getJSONFieldName(input any, field string) string {
	t := reflect.TypeOf(input)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return field
	}

	if f, ok := t.FieldByName(field); ok {
		jsonTag := f.Tag.Get("json")
		name := strings.Split(jsonTag, ",")[0]
		if name != "" && name != "-" {
			return name
		}
	}

	return field
}

func buildValidationErrorMap(input any, verr validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)
	for _, e := range verr {
		field := getJSONFieldName(input, e.Field())
		msg := GetValidationMessage(field, e.Tag(), e.Param())
		errorMessages[field] = msg
	}

	return errorMessages
}
