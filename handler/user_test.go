package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/aoringo0606/echo-project/model"
	"github.com/aoringo0606/echo-project/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type FakeUserRepository struct {
	findByIDResult model.User
	findByIDError error

	createResult model.User
	createError error
}

func (r *FakeUserRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	return r.findByIDResult, r.findByIDError
}

func (r *FakeUserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (r *FakeUserRepository) Create(ctx context.Context, name string, age int) (model.User, error) {
	return r.createResult, r.createError
}

func (r *FakeUserRepository) Update(ctx context.Context, id int, name string, age int) (model.User, error) {
	return model.User{}, nil
}

func (r *FakeUserRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		findByIDResult model.User
		findByIDError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			id:   "1",
			findByIDResult: model.User{
				ID:   1,
				Name: "Alice",
				Age:  20,
			},
			findByIDError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Alice","age":20}` + "\n",
		},
		{
			name:           "not found",
			id:             "999",
			findByIDResult: model.User{},
			findByIDError:  repository.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}` + "\n",
		},
		{
			name:           "internal server error",
			id:             "1",
			findByIDResult: model.User{},
			findByIDError:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}` + "\n",
		},
		{
			name:           "invalid id",
			id:             "abc",
			findByIDResult: model.User{},
			findByIDError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid id"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(
				http.MethodGet,
				"/users/"+tt.id,
				nil,
			)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			fakeRepo := &FakeUserRepository{
				findByIDResult: tt.findByIDResult,
				findByIDError:  tt.findByIDError,
			}

			h := NewUserHandler(fakeRepo)

			err := h.GetUser(c)
			if err != nil {
				t.Fatal(err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatus,
					rec.Code,
				)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf(
					"expected body %q, got %q",
					tt.expectedBody,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string
		body string
		createResult model.User
		createError error
		expectedStatus int
		expectedBody string
	}{
		{
			name: "Success",
			body: `{"name":"Alice","age":20}`,
			createResult: model.User{
				ID: 1,
				Name: "Alice",
				Age: 20,
			},
			createError: nil,
			expectedStatus: http.StatusCreated,
			expectedBody: `{"id":1,"name":"Alice","age":20}` + "\n",
		},
		{
			name: "broken JSON",
			body: `{"name":"Alice","age":`,
			createResult: model.User{},
			createError: nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"error":"invalid request"}` + "\n",
		},
		{
			name: "blank name",
			body: `{"name":"","age":20}`,
			createResult: model.User{},
			createError: nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"errors":{"name":"name is required"}}` + "\n",
		},
		{
			name: "negative age",
			body: `{"name":"Alice","age":-1}`,
			createResult: model.User{},
			createError: nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"errors":{"age":"age must be greater than or equal to 0"}}` + "\n",
		},
		{
			name: "internal server error",
			body: `{"name":"Alice","age":20}`,
			createResult: model.User{},
			createError: errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{"error":"internal server error"}` + "\n",
		},
	}
	for _, tt  := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			v := validator.New()

			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

				if name == "-" {
					return ""
				}

				return name
			})

			e.Validator = NewCustomValidator(v)

						req := httptest.NewRequest(
				http.MethodPost,
				"/users",
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				echo.HeaderContentType,
				echo.MIMEApplicationJSON,
			)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			fakeRepo := &FakeUserRepository{
				createResult: tt.createResult,
				createError: tt.createError,
			}

			h := NewUserHandler(fakeRepo)

			err := h.CreateUser(c)
			if err != nil {
				t.Fatal(err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf(
					"exptected status %d, got %d",
					tt.expectedStatus,
					rec.Code,
				)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf(
					"expected body %q, got %q",
					tt.expectedBody,
					rec.Body.String(),
				)
			}
		})
	}
}