package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepo
}

func NewCategoryService(categoryRepo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) GetAllCategories() ([]model.Category, error) {
	return nil, ErrNotImplemented
}
