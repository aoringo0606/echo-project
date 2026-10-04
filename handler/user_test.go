package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aoringo0606/echo-project/model"
	"github.com/aoringo0606/echo-project/repository"
	"github.com/labstack/echo/v4"
)

type FakeUserRepository struct {
	findByIDResult model.User
	findByIdError error
}

func (r *FakeUserRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	return r.findByIDResult, r.findByIdError
}

func (r *FakeUserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (r *FakeUserRepository) Create(ctx context.Context, name string, age int) (model.User, error) {
	return model.User{}, nil
}

func (r *FakeUserRepository) Update(ctx context.Context, id int, name string, age int) (model.User, error) {
	return model.User{}, nil
}

func (r *FakeUserRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func TestGetUserSuccess(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/1",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	fakeRepo := &FakeUserRepository{
		findByIDResult: model.User{
			ID: 1,
			Name: "Alice",
			Age: 20,
		},
		findByIdError: nil,
	}

	h := NewUserHandler(fakeRepo)

	err := h.GetUser(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	expected := `{"id":1,"name":"Alice","age":20}` + "\n"

	if rec.Body.String() != expected {
		t.Errorf(
			"expected body %q, got %q",
			expected,
			rec.Body.String(),
		)
	}
}

func TestGetUserNotFound(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/999",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("999")

	fakeRepo := &FakeUserRepository{
		findByIDResult: model.User{},
		findByIdError: repository.ErrUserNotFound,
	}

	h := NewUserHandler(fakeRepo)

	err := h.GetUser(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}

	expected := `{"error":"user not found"}` + "\n"

	if rec.Body.String() != expected {
		t.Errorf(
			"expected body %q, got %q",
			expected,
			rec.Body.String(),
		)
	}
}

func TestGetUserInternalServerError(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/500",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("500")

	fakeRepo := &FakeUserRepository{
		findByIDResult: model.User{},
		findByIdError: errors.New("database broken"),
	}

	h := NewUserHandler(fakeRepo)

	err := h.GetUser(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}

	expected := `{"error":"internal server error"}` + "\n"

	if rec.Body.String() != expected {
		t.Errorf(
			"expected body %q, got %q",
			expected,
			rec.Body.String(),
		)
	}
}

func TestGetUserInvalidID(t *testing.T) {
		e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/abc",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("abc")

	fakeRepo := &FakeUserRepository{
		findByIDResult: model.User{},
		findByIdError: errors.New("database error"),
	}

	h := NewUserHandler(fakeRepo)

	err := h.GetUser(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	expected := `{"error":"invalid id"}` + "\n"

	if rec.Body.String() != expected {
		t.Errorf(
			"expected body %q, got %q",
			expected,
			rec.Body.String(),
		)
	}
}