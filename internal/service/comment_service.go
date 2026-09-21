package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"strings"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
}

func NewCommentService(commentRepo *repository.CommentRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(postID, userID int, content string) error {
	comment := &model.Comment{
		PostId:  postID,
		UserID:  userID,
		Content: strings.TrimSpace(content),
	}
	return s.commentRepo.InsertComment(context.Background(), comment)
}

func (s *CommentService) GetCommentsByPost(postID int) ([]model.Comment, error) {
	comments, err := s.commentRepo.ListCommentsByPost(context.Background(), int64(postID))
	if err != nil {
		return nil, err
	}
	if err := s.commentRepo.LoadLikes(context.Background(), comments); err != nil {
		return nil, err
	}
	return toCommentValues(comments), nil
}

func toCommentValues(comments []*model.Comment) []model.Comment {
	values := make([]model.Comment, 0, len(comments))
	for _, comment := range comments {
		values = append(values, *comment)
	}
	return values
}
