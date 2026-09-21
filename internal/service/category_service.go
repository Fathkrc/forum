package service

import "cmd/forum/internal/repository"

type CategoryService struct {
	categoryRepo *repository.CategoryRepo
}

func NewCategoryService(categoryRepo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}
