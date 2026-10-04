package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/aoringo0606/echo-project/model"
	"github.com/aoringo0606/echo-project/repository"
	"github.com/labstack/echo/v4"
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
	GetAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, name string, age int) (model.User, error)
	Update(ctx context.Context, id int, name string, age int) (model.User, error)
	Delete(ctx context.Context, id int) error
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
	users, err := h.repo.GetAll(c.Request().Context())
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	user, err := h.repo.FindByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return errorResponse(c, http.StatusNotFound, "user not found")
		}
		return errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, user)

}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&req); err != nil {
		return handleValidationError(c, err)
	}
	user, err := h.repo.Create(c.Request().Context(), req.Name, req.Age)
	if err != nil {
		return errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
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
	if err := c.Validate(&req); err != nil {
		return handleValidationError(c, err)
	}
	user, err := h.repo.Update(c.Request().Context(), id, req.Name, req.Age)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return errorResponse(c, http.StatusNotFound, "user not found")
		}
		return errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errorResponse(c, http.StatusBadRequest, "invalid id")
	}
	err = h.repo.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return errorResponse(c, http.StatusNotFound, "user not found")
		}
		return errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	return c.NoContent(http.StatusNoContent)
}
