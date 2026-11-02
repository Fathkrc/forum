package repository

import (
	"cmd/forum/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (r *UserRepo) InsertUser(ctx context.Context, user *model.User) error {
	query := `
				INSERT INTO users (email, username, password, created_at)
				VALUES(?,?,?,?)
			`
	result, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)

	return nil

}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
				SELECT id, username, password, created_at 
				FROM users
				WHERE email = ?
			`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx,
		query,
		email).Scan(&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.Email = email
	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
				SELECT email, username
				FROM users
				WHERE id = ?
			`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Email, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	user.ID = int(id)
	return user, nil
}

func (r *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT 1 FROM users WHERE email = ?`,
		email,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT 1 FROM users WHERE username = ?`,
		username,
	).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
