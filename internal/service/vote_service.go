package service

import (
	"cmd/forum/internal/repository"
	"context"
)

type VoteService struct {
	postRepo    *repository.PostRepo
	commentRepo *repository.CommentRepo
}

func NewVoteService(postRepo *repository.PostRepo, commentRepo *repository.CommentRepo) *VoteService {
	return &VoteService{postRepo: postRepo, commentRepo: commentRepo}
}

func (s *VoteService) VotePost(userID, postID int, isLike bool) error {
	return s.postRepo.UpsertPostVote(context.Background(), int64(userID), int64(postID), voteValue(isLike))
}

func (s *VoteService) VoteComment(userID, commentID int, isLike bool) error {
	return s.commentRepo.UpsertCommentVote(context.Background(), int64(userID), int64(commentID), voteValue(isLike))
}

func voteValue(isLike bool) int {
	if isLike {
		return 1
	}
	return -1
}
