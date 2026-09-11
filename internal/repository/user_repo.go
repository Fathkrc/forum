package repository

import (
	"cmd/forum/internal/model"
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *UserRepo) InsertUser(ctx context.Context, user *model.User) error {
	query := `
				INSERT INTO users (email, username, password, create_at)
				VALUES(?,?,?,?)
			`
	result, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}
	user.ID = int(id)

	return nil

}
