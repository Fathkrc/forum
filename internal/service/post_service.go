package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
)

type PostService struct {
	postRepo *repository.PostRepo
}

func NewPostService(postRepo *repository.PostRepo) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(ctx context.Context, userID int, title, content string, categories []model.Category) error {
	post := &model.Post{
		UserID:     userID,
		Title:      title,
		Content:    content,
		Categories: categories}
	return s.postRepo.InsertPost(ctx, post)
}
func (s *PostService) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	return s.postRepo.GetPostByID(ctx, int64(id))
}
func (s *PostService) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	return s.postRepo.ListAllPosts(ctx)
}
func (s *PostService) GetPostsByCategory(ctx context.Context, categoryID int) ([]*model.Post, error) {
	return s.postRepo.ListPostsByCategoryID(ctx, int64(categoryID))
}
func (s *PostService) GetPostsByUser(ctx context.Context, userID int) ([]*model.Post, error) {
	return s.postRepo.ListPostsByUser(ctx, int64(userID))
}
func (s *PostService) GetPostByTerm(ctx context.Context, searchTerm string) ([]*model.Post, error) {
	return s.postRepo.SearchPosts(ctx, searchTerm)
}

//GetLikedPostsByUser(userID int) ([]model.Post, error)
