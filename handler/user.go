package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/aoringo0606/echo-project/model"
)

var users []model.User
var nextID = 1

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0"`
}

func GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK,users)
}

func GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	for _, user := range users {
		if user.ID == id {
			return c.JSON(http.StatusOK, user)
		}
	}
	return errorResponse(c, http.StatusNotFound, "user not found")
}

func CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err!= nil {
		return handleValidationError(c, err)
	}
	user := model.User{
		Name: req.Name,
		Age: req.Age,
	}
	user.ID = nextID
	nextID++
	users = append(users, user)
	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err!= nil {
		return handleValidationError(c, err)
	}
	for i, user := range users {
		if user.ID == id {
			users[i].Name = req.Name
			users[i].Age = req.Age
			return c.JSON(http.StatusOK, users[i])
		}
	}
	return errorResponse(c, http.StatusNotFound, "user not found")
}

func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return errorResponse(c, http.StatusNotFound, "user not found")
}