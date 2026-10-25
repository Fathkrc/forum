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
	voteRepo    *repository.VoteRepo
}

func NewCommentService(commentRepo *repository.CommentRepo, voteRepo *repository.VoteRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo,
		voteRepo: voteRepo}
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

	comments, err := s.commentRepo.ListCommentsByPost(ctx, int64(postID))
	if err != nil {
		return nil, err
	}
	for i := range comments {
		comments[i].Likes, comments[i].Dislikes, err =
			s.voteRepo.GetCommentVoteCounts(ctx, int64(comments[i].ID))
		if err != nil {
			return nil, err
		}
	}
	return comments, nil
}
func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID int64) error {
	comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}

	if int64(comment.UserID) != userID {
		return ErrUnauthorized
	}

	return s.commentRepo.DeleteComment(ctx, commentID)
}
