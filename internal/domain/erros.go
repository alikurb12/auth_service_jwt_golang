package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExsists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrNotVerified        = errors.New("email not verified")
	ErrAlreadyVerified    = errors.New("email already verified")
	ErrForbidden          = errors.New("forbidden")
)
