package repository

import(
	"database/sql"
	"errors"

	"github.com/aoringo0606/echo-project/model"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{
		db: db,
	}
}

func (r *SQLiteUserRepository) GetAll() ([]model.User, error) {
	rows, err := r.db.Query("SELECT id, name, age FROM users")
	if err != nil {
		return []model.User{}, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Age,
		); err != nil {
			return []model.User{}, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return []model.User{}, err
	}
	return users, nil
}

func (r *SQLiteUserRepository) FindByID(id int) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		"SELECT id, name, age FROM users WHERE id = ?",
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) Create(name string, age int) (model.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, age) VALUES (? ?)",
		name,
		age,
	)
	if err != nil {
		return model.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, err
	}
	user := model.User{
		ID: int(id),
		Name: name,
		Age: age,
	}
	return user, nil
}

