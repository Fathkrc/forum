package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
)

type AuthService struct {
	userRepo *repository.UserRepo
}

func NewAuthService(userRepo *repository.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(email, username, password, confirmPassword string) error {
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)

	if email == "" || username == "" || password == "" {
		return errors.New("all fields are required")
	}
	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	exists, err := s.userRepo.EmailExists(context.Background(), email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email is already registered")
	}

	exists, err = s.userRepo.UsernameExists(context.Background(), username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username is already taken")
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	return s.userRepo.InsertUser(context.Background(), user)
}
