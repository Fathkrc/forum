package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
)

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
		VoteService:     NewVoteService(repo.PostRepo, repo.VoteRepo, repo.CommentRepo),
	}
}

func (s *Service) Register(email, username, password, confirmPassword string) error {
	return s.UserService.Register(email, username, password, confirmPassword)
}

func (s *Service) Login(email, password string) (*model.Session, error) {
	return s.UserService.Login(email, password)
}

func (s *Service) Logout(token string) error {
	return s.UserService.Logout(token)
}

func (s *Service) GetSessionByToken(token string) (*model.Session, error) {
	return s.UserService.GetSessionByToken(token)
}

func (s *Service) GetUserByID(id int) (*model.User, error) {
	return s.UserService.GetUserByID(id)
}

func (s *Service) GetAllCategories() ([]model.Category, error) {
	return s.CategoryService.GetAllCategories()
}

func (s *Service) CreateComment(postID, userID int, content string) error {
	return s.CommentService.CreateComment(postID, userID, content)
}

func (s *Service) GetCommentsByPost(postID int) ([]model.Comment, error) {
	return s.CommentService.GetCommentsByPost(postID)
}

func (s *Service) CreatePost(ctx context.Context, userID int, title, content string, categoryIDs []int) error {
	return s.PostService.CreatePost(ctx, userID, title, content, categoryIDs)
}

func (s *Service) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	return s.PostService.GetPostByID(ctx, id)
}

func (s *Service) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	return s.PostService.GetAllPosts(ctx)
}

func (s *Service) GetPostsByCategory(ctx context.Context, categoryID int) ([]model.Post, error) {
	return s.PostService.GetPostsByCategory(ctx, categoryID)
}

func (s *Service) GetPostsByUser(ctx context.Context, userID int) ([]model.Post, error) {
	return s.PostService.GetPostsByUser(ctx, userID)
}

func (s *Service) GetLikedPostsByUser(ctx context.Context, userID int) ([]model.Post, error) {
	return s.PostService.GetLikedPostsByUser(ctx, userID)
}

func (s *Service) SearchPosts(ctx context.Context, query string) ([]model.Post, error) {
	return s.PostService.SearchPosts(ctx, query)
}

func (s *Service) VotePost(ctx context.Context, userID, postID int, isLike bool) error {
	return s.VoteService.VotePost(ctx, userID, postID, isLike)
}

func (s *Service) VoteComment(userID, commentID int, isLike bool) error {
	return s.VoteService.VoteComment(userID, commentID, isLike)
}
