package main

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(loggerMiddleware)

	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	e.Validator = &CustomValidator{
		validator: v,
	}

	userGroup := e.Group("/users")

	userGroup.GET("", getUsers)
	userGroup.GET("/:id", getUser)

	protectedUsers := userGroup.Group("", authMiddleware)

	protectedUsers.POST("", createUser)
	protectedUsers.PUT("/:id", updateUser)
	protectedUsers.DELETE("/:id", deleteUser)

	e.Start(":8080")
}