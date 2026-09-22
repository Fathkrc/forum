package service

import "cmd/forum/internal/repository"

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}
