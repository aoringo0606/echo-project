package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/aoringo0606/echo-project/model"
	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dir := t.TempDir()

	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age INTEGER NOT NULL
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestSQLiteUserRepositoryCRUD(t *testing.T) {
	db := newTestDB(t)
	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Alice", 20)
	if err != nil {
		t.Fatal(err)
	}

	expected := model.User{
		ID: 1,
		Name: "Alice",
		Age: 20,
	}

	if created != expected {
		t.Errorf(
			"expected user %+v, got %+v",
			expected,
			created,
		)
	}

	found, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if found != created {
		t.Errorf(
			"expected user %+v, got %+v",
			created,
			found,
		)
	}

	updated, err := repo.Update(ctx, created.ID, "Alice Updated", 21)
	if err != nil {
		t.Fatal(err)
	}

	expected = model.User{
		ID: 1,
		Name: "Alice Updated",
		Age: 21,
	}

	if updated != expected {
		t.Errorf(
			"expected user %+v, got %+v",
			expected,
			updated,
		)
	}

	found, err = repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if found != updated {
		t.Errorf(
			"expected user %+v, got %+v",
			updated,
			found,
		)
	}

	users, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf(
			"expected users size %d, got %d",
			1,
			len(users),
		)
	}
	if users[0] != updated {
		t.Errorf(
			"expected user %+v, got %+v",
			updated,
			users[0],
		)
	}

	err = repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.FindByID(ctx, created.ID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}

func TestSQLiteUserRepositoryUpdateNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	_, err := repo.Update(ctx, 999, "Alice", 20)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}

func TestSQLiteUserRepositoryDeleteNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 1)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}