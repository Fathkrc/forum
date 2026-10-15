package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"errors"
	"strings"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
}

func NewCommentService(commentRepo *repository.CommentRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(ctx context.Context, postID, userID int, content string) error {
	if strings.TrimSpace(content) == "" {
		return errors.New("Invalid Comment")
	}
	comment := &model.Comment{
		PostId:  postID,
		UserID:  userID,
		Content: content}
	return s.commentRepo.InsertComment(ctx, comment)
}

func (s *CommentService) GetCommentsByPost(ctx context.Context, postID int) ([]*model.Comment, error) {
	return s.commentRepo.ListCommentsByPost(ctx, int64(postID))
}
