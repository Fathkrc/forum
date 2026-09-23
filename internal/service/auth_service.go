package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
)

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *AuthService) Register(email, username, password, confirmPassword string) error {
	return ErrNotImplemented
}

func (s *AuthService) Login(email, password string) (*model.Session, error) {
	return nil, ErrNotImplemented
}

func (s *AuthService) Logout(token string) error {
	return ErrNotImplemented
}

func (s *AuthService) GetSessionByToken(token string) (*model.Session, error) {
	return nil, ErrNotImplemented
}

func (s *AuthService) GetUserByID(id int) (*model.User, error) {
	return nil, ErrNotImplemented
}
