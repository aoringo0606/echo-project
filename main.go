package main

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/aoringo0606/echo-project/handler"
)

func main() {
	e := echo.New()

	e.Use(handler.LoggerMiddleware)

	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	e.Validator = handler.NewCustomValidator(v)

	userGroup := e.Group("/users")

	userGroup.GET("", handler.GetUsers)
	userGroup.GET("/:id", handler.GetUser)

	protectedUsers := userGroup.Group("", handler.AuthMiddleware)

	protectedUsers.POST("", handler.CreateUser)
	protectedUsers.PUT("/:id", handler.UpdateUser)
	protectedUsers.DELETE("/:id", handler.DeleteUser)

	e.Start(":8080")
}