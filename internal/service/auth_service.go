package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

const sessionDuration = 24 * time.Hour

func NewAuthService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *AuthService) Register(email, username, password, confirmPassword string) error {
	email = normalizeEmail(email)
	username = strings.TrimSpace(username)

	if err := validateRegistration(email, username, password, confirmPassword); err != nil {
		return err
	}

	ctx := context.Background()

	emailExists, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return err
	}
	if emailExists {
		return ErrEmailExists
	}

	usernameExists, err := s.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return err
	}
	if usernameExists {
		return ErrUsernameExists
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}

	return s.userRepo.InsertUser(ctx, user)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateRegistration(email, username, password, confirmPassword string) error {
	if email == "" {
		return ErrEmailRequired
	}

	address, err := mail.ParseAddress(email)
	if err != nil || !strings.EqualFold(address.Address, email) {
		return ErrInvalidEmail
	}

	if username == "" {
		return ErrUsernameRequired
	}

	if strings.TrimSpace(password) == "" {
		return ErrPasswordRequired
	}

	if confirmPassword == "" {
		return ErrPasswordConfirmationRequired
	}

	if password != confirmPassword {
		return ErrPasswordsDoNotMatch
	}

	return nil
}

func newSessionToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func (s *AuthService) Login(email, password string) (*model.Session, error) {
	ctx := context.Background()
	user, err := s.userRepo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	matches, err := VerifyPassword(password, user.PasswordHash)
	if err != nil || !matches {
		return nil, ErrInvalidCredentials
	}

	token, err := newSessionToken()
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(sessionDuration),
	}
	if err := s.sessionRepo.InsertSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.DeleteSession(context.Background(), token)
}

func (s *AuthService) GetSessionByToken(token string) (*model.Session, error) {
	ctx := context.Background()
	session, err := s.sessionRepo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !time.Now().UTC().Before(session.ExpiresAt) {
		_ = s.sessionRepo.DeleteSession(ctx, token)
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *AuthService) GetUserByID(id int) (*model.User, error) {
	return s.userRepo.GetUserByID(context.Background(), int64(id))
}
