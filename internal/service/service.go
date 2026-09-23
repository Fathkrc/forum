package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
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
		VoteService:     NewVoteService(repo.PostRepo, repo.CommentRepo),
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

func (s *Service) CreatePost(userID int, title, content string, categoryIDs []int) error {
	return s.PostService.CreatePost(userID, title, content, categoryIDs)
}

func (s *Service) GetPostByID(id int) (*model.Post, error) {
	return s.PostService.GetPostByID(id)
}

func (s *Service) GetAllPosts() ([]model.Post, error) {
	return s.PostService.GetAllPosts()
}

func (s *Service) GetPostsByCategory(categoryID int) ([]model.Post, error) {
	return s.PostService.GetPostsByCategory(categoryID)
}

func (s *Service) GetPostsByUser(userID int) ([]model.Post, error) {
	return s.PostService.GetPostsByUser(userID)
}

func (s *Service) GetLikedPostsByUser(userID int) ([]model.Post, error) {
	return s.PostService.GetLikedPostsByUser(userID)
}

func (s *Service) SearchPosts(query string) ([]model.Post, error) {
	return s.PostService.SearchPosts(query)
}

func (s *Service) VotePost(userID, postID int, isLike bool) error {
	return s.VoteService.VotePost(userID, postID, isLike)
}

func (s *Service) VoteComment(userID, commentID int, isLike bool) error {
	return s.VoteService.VoteComment(userID, commentID, isLike)
}
