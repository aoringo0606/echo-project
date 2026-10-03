package repository

import(
	"github.com/aoringo0606/echo-project/model"
)

type UserRepository struct {
	users  []model.User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  []model.User{},
		nextID: 1,
	}
}

func (r *UserRepository) GetAll() []model.User {
	return r.users
}

func (r *UserRepository) Create(name string, age int) model.User {
	user := model.User{
		ID: r.nextID,
		Name: name,
		Age: age,
	}

	r.nextID++
	r.users = append(r.users, user)

	return user
}

func (r *UserRepository) FindByID(id int) (model.User, bool) {
	for _, user := range r.users {
		if user.ID == id {
			return user, true
		}
	}
	return model.User{}, false
}

func (r *UserRepository) Update(id int, name string, age int) (model.User, bool) {
	for i, user := range r.users {
		if user.ID == id {
			r.users[i].Name = name
			r.users[i].Age = age
		}
	}
	return model.User{}, false
}

func (r *UserRepository) Delete(id int) bool {
	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return true
		}
	}
	return false
}