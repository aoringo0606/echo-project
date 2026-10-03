package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/aoringo0606/echo-project/model"
)

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0"`
}

type UserRepository interface {
	GetAll() []model.User
	FindByID(id int) (model.User, bool)
	Create(name string, age int) model.User
	Update(id int, name string, age int) (model.User, bool)
	Delete(id int) bool
}

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users := h.repo.GetAll()
	return c.JSON(http.StatusOK, users)
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
	user, ok := h.repo.Update(id, req.Name, req.Age)
	if !ok {
		return errorResponse(c, http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	ok := h.repo.Delete(id)
	if !ok {
		return errorResponse(c, http.StatusNotFound, "user not found")
	}
	return c.NoContent(http.StatusNoContent)
}