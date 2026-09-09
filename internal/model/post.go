package model

import (
	"time"
)

type Post struct {
	ID             int
	UserID         int
	Title          string
	Content        string
	CreatedAt      time.Time
	AuthorUsername string
	Categories     []Category
	Likes          int
	Dislikes       int
}
