package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepo
}

func NewCategoryService(categoryRepo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	err := s.categoryRepo.InsertCategory(ctx, &category)
	if err != nil {
		return category, nil
	}
	return model.Category{}, err
}
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	categoryList, err := s.categoryRepo.ListMultiCategory(ctx)
	if err != nil {
		return categoryList, nil
	}
	return nil, err
}
