package repository

import (
	"cmd/forum/internal/model"
	"context"
	"database/sql"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) InsertSession(ctx context.Context, session *model.Session) error {
	query := `INSERT INTO sessions (uuid, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, session.Token, session.UserID, session.ExpiresAt)
	return err
}

func (r *SessionRepo) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	query := `SELECT uuid, user_id, expires_at FROM sessions WHERE uuid = ?`
	session := &model.Session{}
	if err := r.db.QueryRowContext(ctx, query, token).
		Scan(&session.Token, &session.UserID, &session.ExpiresAt); err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE uuid = ?`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
