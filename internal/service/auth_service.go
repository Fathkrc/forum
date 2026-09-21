package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
)

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
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
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	return s.userRepo.InsertUser(context.Background(), user)
}

func (s *AuthService) Login(email, password string) (*model.Session, error) {
	user, err := s.userRepo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil || !match {
		return nil, errors.New("invalid email or password")
	}

	token, err := newSessionToken()
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.sessionRepo.InsertSession(context.Background(), session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.DeleteSession(context.Background(), token)
}

func (s *AuthService) GetSessionByToken(token string) (*model.Session, error) {
	session, err := s.sessionRepo.GetSessionByToken(context.Background(), token)
	if err != nil {
		return nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}
	return session, nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
