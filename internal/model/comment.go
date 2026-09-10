package model

type Comment struct {
	ID           int
	PostId       int
	UserID       int
	AuthorUserID int
	Content      string
	Likes        int
	Dislikes     int
}
