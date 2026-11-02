package service

import "errors"

var (
	ErrEmailRequired                = errors.New("email is required")
	ErrInvalidEmail                 = errors.New("please enter a valid email address")
	ErrUsernameRequired             = errors.New("username is required")
	ErrPasswordRequired             = errors.New("password is required")
	ErrPasswordConfirmationRequired = errors.New("password confirmation is required")
	ErrPasswordsDoNotMatch          = errors.New("passwords do not match")
	ErrEmailExists                  = errors.New("an account with this email already exists")
	ErrUsernameExists               = errors.New("that username is already taken")
	ErrInvalidUsername              = errors.New("user name must be between 4 and 12 characters")
	ErrInvalidCredentials           = errors.New("invalid email or password")
	ErrSessionExpired               = errors.New("session expired")
	ErrNotImplemented               = errors.New("not implemented")
	ErrUnauthorized                 = errors.New("unauthorized")
	ErrUserNotFound                 = errors.New("user not found")
	ErrSessionNotFound              = errors.New("session not found")
)
