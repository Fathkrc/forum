package repository

import "database/sql"

type Repository struct {
	catRepo     *CategoryRepo
	commentRepo *CommentRepo
	userRepo    *UserRepo
	postRepo    *PostRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		catRepo:     NewCategoryRepo(db),
		commentRepo: NewCommentRepo(db),
		userRepo:    NewUserRepo(db),
		postRepo:    NewPostRepo(db),
	}
}
