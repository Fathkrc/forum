package model

import "time"

type Session struct {
	ID        int
	Token     string // token is UUID
	UserID    int
	ExpiresAt time.Time
	CreatedAt time.Time
}
