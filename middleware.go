package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println(c.Request().Method + " " + c.Request().URL.Path)

		start := time.Now()

		err := next(c)

		elapsed := time.Since(start)

		fmt.Println("elapsed: ", elapsed)

		return err
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 {
			return errorResponse(c, http.StatusUnauthorized, "unauthorized")
		}
		if parts[0] != "Bearer" {
			return errorResponse(c, http.StatusUnauthorized, "unauthorized")
		}
		if parts[1] != "secret-token" {
			return errorResponse(c, http.StatusUnauthorized, "unauthorized")
		}

		return next(c)
	}
}