package main

import (
	"reflect"
	"strings"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
	"github.com/aoringo0606/echo-project/handler"
	"github.com/aoringo0606/echo-project/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(handler.LoggerMiddleware)

	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewSQLiteUserRepository(db)
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