package model

type Comment struct {
	ID             int
	PostId         int
	UserID         int
	AuthorUsername string
	Content        string
	Likes          int
	Dislikes       int
}
