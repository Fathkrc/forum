package model

import "time"

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
}

type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	CreatedAt    time.Time
	Username     string
	LikeCount    int
	DislikeCount int
	CommentCount int
	Categories   []Category
}

type Comment struct {
	ID           int
	PostID       int
	UserID       int
	Content      string
	CreatedAt    time.Time
	Username     string
	LikeCount    int
	DislikeCount int
}

type Category struct {
	ID   int
	Name string
}
