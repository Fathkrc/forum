package service

import "cmd/forum/internal/repository"

type PostService struct {
	postRepo *repository.PostRepo
}

func NewPostService(postRepo *repository.PostRepo) *PostService {
	return &PostService{postRepo: postRepo}
}
