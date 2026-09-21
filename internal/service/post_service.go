package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"strings"
)

type PostService struct {
	postRepo     *repository.PostRepo
	categoryRepo *repository.CategoryRepo
}

func NewPostService(postRepo *repository.PostRepo, categoryRepo *repository.CategoryRepo) *PostService {
	return &PostService{postRepo: postRepo, categoryRepo: categoryRepo}
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
func (s *PostService) GetPostByID(id int) (*model.Post, error) {
	post, err := s.postRepo.GetPostByID(context.Background(), int64(id))
	if err != nil {
		return nil, err
	}
	if err := s.enrich([]*model.Post{post}); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetAllPosts() ([]model.Post, error) {
	posts, err := s.postRepo.ListAllPosts(context.Background())
	if err != nil {
		return nil, err
	}
	if err := s.enrich(posts); err != nil {
		return nil, err
	}
	return toPostValues(posts), nil
}

func (s *PostService) GetPostsByCategory(categoryID int) ([]model.Post, error) {
	posts, err := s.postRepo.ListPostsByCategoryID(context.Background(), int64(categoryID))
	if err != nil {
		return nil, err
	}
	if err := s.enrich(posts); err != nil {
		return nil, err
	}
	return toPostValues(posts), nil
}

func (s *PostService) GetPostsByUser(userID int) ([]model.Post, error) {
	posts, err := s.postRepo.ListPostsByUser(context.Background(), int64(userID))
	if err != nil {
		return nil, err
	}
	if err := s.enrich(posts); err != nil {
		return nil, err
	}
	return toPostValues(posts), nil
}

func (s *PostService) GetPostByTerm(ctx context.Context, searchTerm string) ([]*model.Post, error) {
	return s.postRepo.SearchPosts(ctx, searchTerm)
}

func (s *PostService) enrich(posts []*model.Post) error {
	if err := s.postRepo.LoadDetails(context.Background(), posts); err != nil {
		return err
	}
	for _, post := range posts {
		if err := s.loadCategories(post); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) loadCategories(post *model.Post) error {
	categories, err := s.categoryRepo.ListCategoriesByPost(context.Background(), int64(post.ID))
	if err != nil {
		return err
	}
	for _, category := range categories {
		post.Categories = append(post.Categories, *category)
	}
	return nil
}

func toPostValues(posts []*model.Post) []model.Post {
	values := make([]model.Post, 0, len(posts))
	for _, post := range posts {
		values = append(values, *post)
	}
	return values
}
