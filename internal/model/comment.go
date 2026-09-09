package model

type Comment struct {
	ID           int
	PostId       int
	UserID       int
	Categories   []Category
	AuthorUserID int
	Likes        int
	Dislikes     int
}
