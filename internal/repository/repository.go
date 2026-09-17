package repository

import "database/sql"

type Repositoy struct {
	catRepo     *CategoryRepo
	commentRepo *CommentRepo
	userRepo    *UserRepo
	postRepo    *PostRepo
}

func NewRepository(db *sql.DB) *Repositoy {
	return &Repositoy{
		catRepo:     NewCategoryRepo(db),
		commentRepo: NewCommentRepo(db),
		userRepo:    NewUserRepo(db),
		postRepo:    NewPostRepo(db),
	}
}
