package service

import "cmd/forum/internal/repository"

type Service struct {
	UserService     *AuthService
	CategoryService *CategoryService
	CommentService  *CommentService
	PostService     *PostService
	VoteService     *VoteService
}

func New(repo *repository.Repository) *Service {
	return &Service{
		UserService:     NewAuthService(repo.UserRepo, repo.SessionRepo),
		CategoryService: NewCategoryService(repo.CatRepo),
		CommentService:  NewCommentService(repo.CommentRepo),
		PostService:     NewPostService(repo.PostRepo),
		VoteService:     NewVoteService(repo.PostRepo, repo.CommentRepo),
	}
}
