package model

type Comment struct {
	ID           int
	PostId       int
	UserID       int
	AuthorUserID int
	Likes        int
	Dislikes     int
}
