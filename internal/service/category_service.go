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

func (s *CategoryService) GetAllCategories() ([]model.Category, error) {
	categories, err := s.categoryRepo.ListMultiCategory(context.Background())
	if err != nil {
		return nil, err
	}

	values := make([]model.Category, 0, len(categories))
	for _, category := range categories {
		values = append(values, *category)
	}
	return values, nil
}
