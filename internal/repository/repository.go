package repository

import "database/sql"

type Repository struct {
	CatRepo     *CategoryRepo
	CommentRepo *CommentRepo
	UserRepo    *UserRepo
	PostRepo    *PostRepo
	SessionRepo *SessionRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		CatRepo:     NewCategoryRepo(db),
		CommentRepo: NewCommentRepo(db),
		UserRepo:    NewUserRepo(db),
		PostRepo:    NewPostRepo(db),
		SessionRepo: NewSessionRepo(db),
	}
}
