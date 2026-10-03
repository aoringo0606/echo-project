package handler

import "github.com/labstack/echo/v4"

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

type ErrorMessageResponse struct {
	Error string `json:"error"`
}

func errorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrorMessageResponse{
		Error: message,
	})
}
