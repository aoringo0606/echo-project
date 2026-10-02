package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator(v *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func handleValidationError(c echo.Context, err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errorResponse(c, http.StatusInternalServerError, "validation error")
	}

	errors := make(map[string]string)

	for _, validationErr := range validationErrors {
		field := validationErr.Field()

		switch validationErr.Tag() {
		case "required":
			errors[field] = field + " is required"
		case "gte":
			errors[field] = field + " must be greater than or equal to " +
				validationErr.Param()
		default:
			errors[field] = "invalid value"
		}
	}
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Errors: errors,
	})
}