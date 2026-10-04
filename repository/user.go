package repository

import (
	"context"
	"errors"

	"github.com/aoringo0606/echo-project/model"
)

type UserRepository struct {
	users  []model.User
	nextID int
}

var ErrUserNotFound = errors.New("user not found")

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  []model.User{},
		nextID: 1,
	}
}

func (r *UserRepository) GetAll(_ context.Context) ([]model.User, error) {
	return r.users, nil
}

func (r *UserRepository) Create(_ context.Context, name string, age int) (model.User, error) {
	user := model.User{
		ID:   r.nextID,
		Name: name,
		Age:  age,
	}

	r.nextID++
	r.users = append(r.users, user)

	return user, nil
}

func (r *UserRepository) FindByID(_ context.Context, id int) (model.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return model.User{}, ErrUserNotFound
}

func (r *UserRepository) Update(_ context.Context, id int, name string, age int) (model.User, error) {
	for i, user := range r.users {
		if user.ID == id {
			r.users[i].Name = name
			r.users[i].Age = age
			return r.users[i], nil
		}
	}
	return model.User{}, ErrUserNotFound
}

func (r *UserRepository) Delete(_ context.Context, id int) error {
	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}
