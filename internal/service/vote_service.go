package service

import "cmd/forum/internal/repository"

type VoteService struct {
	postRepo    *repository.PostRepo
	commentRepo *repository.CommentRepo
}

func NewVoteService(postRepo *repository.PostRepo, commentRepo *repository.CommentRepo) *VoteService {
	return &VoteService{postRepo: postRepo, commentRepo: commentRepo}
}

func (s *VoteService) VotePost(userID, postID int, isLike bool) error {
	return ErrNotImplemented
}

func (s *VoteService) VoteComment(userID, commentID int, isLike bool) error {
	return ErrNotImplemented
}
