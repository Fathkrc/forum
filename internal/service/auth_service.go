package service

import (
	"cmd/forum/internal/model"
	"cmd/forum/internal/repository"
	"net/mail"
	"strings"
)

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *AuthService) Register(email, username, password, confirmPassword string) error {
	email = normalizeEmail(email)
	username = strings.TrimSpace(username)

	if err := validateRegistration(email, username, password, confirmPassword); err != nil {
		return err
	}

	return ErrNotImplemented
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
