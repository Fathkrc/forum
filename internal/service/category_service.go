package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"errors"
	"strings"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepo
}

func NewCategoryService(categoryRepo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category model.Category) (*model.Category, error) {
	category.Name = strings.TrimSpace(category.Name)

	if category.Name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	if err := s.categoryRepo.InsertCategory(ctx, &category); err != nil {
		return nil, err
	}

	return &category, nil
}
func (s *CategoryService) GetCategoryByID(
	ctx context.Context,
	id int,
) (*model.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid category ID")
	}

	return s.categoryRepo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) UpdateCategory(
	ctx context.Context,
	category model.Category,
) error {
	category.Name = strings.TrimSpace(category.Name)

	if category.ID <= 0 {
		return errors.New("invalid category ID")
	}
	if category.Name == "" {
		return errors.New("category name cannot be empty")
	}

	return s.categoryRepo.UpdateCategory(ctx, &category)
}

func (s *CategoryService) DeleteCategory(
	ctx context.Context,
	id int,
) error {
	if id <= 0 {
		return errors.New("invalid category ID")
	}

	return s.categoryRepo.DeleteCategory(ctx, id)
}
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return s.categoryRepo.ListMultiCategory(ctx)
}
