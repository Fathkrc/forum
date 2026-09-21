package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
)

type AuthService struct {
	userRepo *repository.UserRepo
}

func NewAuthService(userRepo *repository.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, email, username, password, confirmPassword string) error {

	user := &model.User{}
	return s.userRepo.InsertUser(ctx, user)
}
