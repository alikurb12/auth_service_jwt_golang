package usecase

import (
	"context"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.AuthResponse, error)
	Login(ctx context.Context, input domain.LoginInput) (*domain.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerifiction(ctx context.Context, email string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
