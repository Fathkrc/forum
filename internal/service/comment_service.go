package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
}

func NewCommentService(commentRepo *repository.CommentRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(postID, userID int, content string) error {
	return ErrNotImplemented
}

func (s *CommentService) GetCommentsByPost(postID int) ([]model.Comment, error) {
	return nil, ErrNotImplemented
}
