package service

import "cmd/forum/internal/repository"

type VoteService struct {
	postRepo    *repository.PostRepo
	commentRepo *repository.CommentRepo
}

func NewVoteService(postRepo *repository.PostRepo, commentRepo *repository.CommentRepo) *VoteService {
	return &VoteService{postRepo: postRepo, commentRepo: commentRepo}
}
