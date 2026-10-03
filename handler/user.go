package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/aoringo0606/echo-project/model"
	"github.com/aoringo0606/echo-project/repository"
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

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	return c.JSON(http.StatusOK,users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	user, ok := h.repo.FindByID(id)
	if !ok {
		return errorResponse(c, http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
	
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err!= nil {
		return handleValidationError(c, err)
	}
	user := h.repo.Create(req.Name, req.Age)
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
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

func (h *UserHandler) DeleteUser(c echo.Context) error {
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