package service

import "cmd/forum/internal/repository"

type ForumService struct {
	UserService     *AuthService
	CategoryService *CategoryService
	CommentService  *CommentService
	PostService     *PostService
	VoteService     *VoteService
}

func New(repo *repository.Repository) *ForumService {
	return &ForumService{
		UserService:     NewAuthService(repo.UserRepo, repo.SessionRepo),
		CategoryService: NewCategoryService(repo.CatRepo),
		CommentService:  NewCommentService(repo.CommentRepo),
		PostService:     NewPostService(repo.PostRepo),
		VoteService:     NewVoteService(repo.PostRepo, repo.CommentRepo),
	}
}
