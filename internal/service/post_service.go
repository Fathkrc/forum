package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepo
}

func NewPostService(postRepo *repository.PostRepo) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(userID int, title, content string, categoryIDs []int) error {
	return ErrNotImplemented
}

func (s *PostService) GetPostByID(id int) (*model.Post, error) {
	return nil, ErrNotImplemented
}

func (s *PostService) GetAllPosts() ([]model.Post, error) {
	return nil, ErrNotImplemented
}

func (s *PostService) GetPostsByCategory(categoryID int) ([]model.Post, error) {
	return nil, ErrNotImplemented
}

func (s *PostService) GetPostsByUser(userID int) ([]model.Post, error) {
	return nil, ErrNotImplemented
}

func (s *PostService) GetLikedPostsByUser(userID int) ([]model.Post, error) {
	return nil, ErrNotImplemented
}

func (s *PostService) SearchPosts(query string) ([]model.Post, error) {
	return nil, ErrNotImplemented
}
