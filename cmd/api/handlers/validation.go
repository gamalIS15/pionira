package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"pionira/common"
	"reflect"
	"strings"
)

func (h *Handler) ValidateBody(c echo.Context, payload interface{}) []*common.ValidationError {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	var errors []*common.ValidationError

	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		reflected := reflect.ValueOf(payload)
		for _, err := range validationErrors {

			field, _ := reflected.Type().FieldByName(err.StructField())
			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(err.StructField())
			}
			condition := err.Tag()
			keyToSpace := strings.Replace(key, "_", " ", -1)
			//errMessage := keyToSpace + " " + condition
			var errMessage string
			switch condition {
			case "required":
				errMessage = keyToSpace + " is required"
			case "email":
				errMessage = keyToSpace + " must be a valid email address"
			case "min":
				errMessage = keyToSpace + " password length must be at least 8 characters"
			case "eqfield":
				errMessage = keyToSpace + " field name must be equal to " + strings.ToLower(err.Param())

			}

			currentValidationError := &common.ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, currentValidationError)
		}
	}

	return errors
}
