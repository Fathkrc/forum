package service

import (
	"cmd/forum/internal/repository"
	"context"
	"errors"
)

type VoteService struct {
	postRepo    *repository.PostRepo
	commentRepo *repository.CommentRepo
	voteRepo    *repository.VoteRepo
}

func NewVoteService(postRepo *repository.PostRepo, voteRepo *repository.VoteRepo, commentRepo *repository.CommentRepo) *VoteService {
	return &VoteService{postRepo: postRepo, commentRepo: commentRepo, voteRepo: voteRepo}
}

func (s *VoteService) VotePost(ctx context.Context, userID, postID int, isLike bool) error {
	if userID <= 0 || postID <= 0 {
		return errors.New("invalid user or post ID")
	}
	vote := -1
	if isLike {
		vote = 1
	}
	return s.voteRepo.VotePost(
		context.Background(),
		int64(userID),
		int64(postID),
		vote,
	)
}

func (s *VoteService) VoteComment(userID, commentID int, isLike bool) error {
	return ErrNotImplemented
}
