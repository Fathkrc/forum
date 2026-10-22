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
		CommentService:  NewCommentService(repo.CommentRepo, repo.VoteRepo),
		PostService:     NewPostService(repo.PostRepo, repo.VoteRepo, repo.CommentRepo),
		VoteService:     NewVoteService(repo.PostRepo, repo.VoteRepo, repo.CommentRepo),
	}
}

func (s *Service) Register(ctx context.Context, email, username, password, confirmPassword string) error {
	return s.UserService.Register(ctx, email, username, password, confirmPassword)
}

func (s *Service) Login(ctx context.Context, email, password string) (*model.Session, error) {
	return s.UserService.Login(ctx, email, password)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.UserService.Logout(ctx, token)
}

func (s *Service) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	return s.UserService.GetSessionByToken(ctx, token)
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return s.UserService.GetUserByID(ctx, id)
}

func (s *Service) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return s.CategoryService.GetAllCategories(ctx)
}
func (s *Service) CreateCategory(
	ctx context.Context,
	category model.Category,
) (*model.Category, error) {
	return s.CategoryService.CreateCategory(ctx, category)
}

func (s *Service) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	return s.CategoryService.GetCategoryByID(ctx, id)
}

func (s *Service) UpdateCategory(ctx context.Context, category model.Category) error {
	return s.CategoryService.UpdateCategory(ctx, category)
}

func (s *Service) DeleteCategory(ctx context.Context, id int) error {
	return s.CategoryService.DeleteCategory(ctx, id)
}

func (s *Service) CreateComment(ctx context.Context, postID, userID int, content string) error {
	return s.CommentService.CreateComment(ctx, postID, userID, content)
}

func (s *Service) GetCommentsByPost(ctx context.Context, postID int) ([]*model.Comment, error) {
	return s.CommentService.GetCommentsByPost(ctx, postID)
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

func (s *Service) VoteComment(ctx context.Context, userID, commentID int, isLike bool) error {
	return s.VoteService.VoteComment(ctx, userID, commentID, isLike)
}

func (s *Service) GetPostVote(ctx context.Context,
	userID int64,
	postID int64,
) (int, error) {
	return s.VoteService.GetPostVote(ctx, userID, postID)
}

func (s *Service) GetCommentVote(
	ctx context.Context,
	userID int64,
	commentID int64,
) (int, error) {
	return s.VoteService.GetCommentVote(ctx, userID, commentID)
}

func (s *Service) DeletePostVote(
	ctx context.Context,
	userID int64,
	postID int64,
) error {
	return s.VoteService.DeletePostVote(ctx, userID, postID)
}

func (s *Service) DeleteCommentVote(ctx context.Context,
	userID int64,
	commentID int64,
) error {
	return s.VoteService.DeleteCommentVote(ctx, userID, commentID)
}

func (s *Service) GetPostVoteCounts(ctx context.Context, postID int64) (int, int, error) {
	return s.VoteService.GetPostVoteCounts(ctx, postID)
}

func (s *Service) GetCommentVoteCounts(ctx context.Context, commentID int64) (int, int, error) {
	return s.VoteService.GetCommentVoteCounts(ctx, commentID)
}

func (s *Service) DeletePost(ctx context.Context, postID int64, userID int64) error {
	return s.PostService.DeletePost(ctx, postID, userID)
}
