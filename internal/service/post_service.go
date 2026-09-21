package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"strings"
)

type PostService struct {
	postRepo *repository.PostRepo
}

func NewPostService(postRepo *repository.PostRepo) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(userID int, title, content string, categoryIDs []int) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	post := &model.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	if err := s.postRepo.InsertPost(context.Background(), post); err != nil {
		return err
	}

	for _, categoryID := range categoryIDs {
		if err := s.postRepo.AddCategoryToPost(context.Background(), int64(post.ID), int64(categoryID)); err != nil {
			return err
		}
	}
	return nil
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
