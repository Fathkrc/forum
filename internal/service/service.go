package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
)

type ForumService struct {
	UserService     *AuthService
	CategoryService *CategoryService
	CommentService  *CommentService
	PostService     *PostService
}

func New(repository *repository.Repository) *ForumService {
	return &ForumService{
		UserService:     NewAuthService(repository.UserRepo, repository.SessionRepo),
		CategoryService: NewCategoryService(repository.CatRepo),
		CommentService:  NewCommentService(repository.CommentRepo),
		PostService:     NewPostService(repository.PostRepo),
	}
}

func (f *ForumService) Register(email, username, password, confirmPassword string) error {
	return f.UserService.Register(email, username, password, confirmPassword)
}

func (f *ForumService) Login(email, password string) (*model.Session, error) {
	return f.UserService.Login(email, password)
}

func (f *ForumService) Logout(token string) error {
	return f.UserService.Logout(token)
}

func (f *ForumService) GetSessionByToken(token string) (*model.Session, error) {
	return f.UserService.GetSessionByToken(token)
}

func (f *ForumService) GetUserByID(id int) (*model.User, error) {
	return f.UserService.GetUserByID(id)
}

func (f *ForumService) GetAllCategories() ([]model.Category, error) {
	return f.CategoryService.GetAllCategories()
}

type Service interface {
	Register(email, username, password, confirmPassword string) error
	Login(email, password string) (*model.Session, error)
	Logout(token string) error //deletes a session by its token

	CreatePost(userID int, title, content string, categoryIDs []int) error
	GetPostByID(id int) (*model.Post, error)
	GetAllPosts() ([]model.Post, error) // This powers the homepage.
	GetPostsByCategory(categoryID int) ([]model.Post, error)
	GetPostsByUser(userID int) ([]model.Post, error)
	GetLikedPostsByUser(userID int) ([]model.Post, error)

	CreateComment(postID, userID int, content string) error
	GetCommentsByPost(postID int) ([]model.Comment, error)

	VotePost(userID, postID int, isLike bool) error
	VoteComment(userID, commentID int, isLike bool) error

	GetAllCategories() ([]model.Category, error)
	GetSessionByToken(token string) (*model.Session, error)

	// Bonus feature (search page & profile page)
	SearchPosts(query string) ([]model.Post, error)
	GetUserByID(id int) (*model.User, error)
}
