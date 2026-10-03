package main

import (
	"reflect"
	"strings"

	"github.com/aoringo0606/echo-project/handler"
	"github.com/aoringo0606/echo-project/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(handler.LoggerMiddleware)

	repo := repository.NewUserRepository()
	userHandler := handler.NewUserHandler(repo)

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

	userGroup.GET("", userHandler.GetUsers)
	userGroup.GET("/:id", userHandler.GetUser)

	protectedUsers := userGroup.Group("", handler.AuthMiddleware)

	protectedUsers.POST("", userHandler.CreateUser)
	protectedUsers.PUT("/:id", userHandler.UpdateUser)
	protectedUsers.DELETE("/:id", userHandler.DeleteUser)

	e.Start(":8080")
}