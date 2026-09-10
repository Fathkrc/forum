package model

import "time"

type Comment struct {
	ID             int
	PostId         int
	UserID         int
	AuthorUsername string
	Content        string
	Likes          int
	Dislikes       int
	CreatedAt      time.Time
}
