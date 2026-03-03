package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/alikurb12/auth_service_jwt_golang/internal/infrastructure/jwt"
	"github.com/alikurb12/auth_service_jwt_golang/internal/repository"
	"github.com/alikurb12/auth_service_jwt_golang/pkg/hasher"
	"github.com/google/uuid"
)

type EmailSender interface {
	SendVerification(toEmail, token string) error
}

type authUseCase struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtService *jwt.JWTService
	email      EmailSender
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtService *jwt.JWTService,
	email EmailSender,
) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
		email:      email,
	}
}

func (uc *authUseCase) Register(ctx context.Context, input domain.RegisterInput) (*domain.AuthResponse, error) {
	_, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.ErrUserAlreadyExsists
	}
	if err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("Register check email error: %w", err)
	}

	passwordHash, err := hasher.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("Register hash password error: %w", err)
	}

	verifyToken, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("Register generate verify token error: %w", err)
	}

	user := &domain.User{
		ID:             uuid.New(),
		Email:          input.Email,
		PasswordHash:   passwordHash,
		Role:           domain.RoleUser,
		IsVerified:     false,
		VerifyToken:    verifyToken,
		VerifyTokenExp: time.Now().Add(24 * time.Hour),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	if err := uc.email.SendVerification(user.Email, verifyToken); err != nil {
		fmt.Printf("[WARN] send verification email to %s: %v\n", user.Email, err)
	}

	return uc.issueTokens(ctx, user)
}

func (uc *authUseCase) issueTokens(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
	access_token, err := uc.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("Issue tokens generate access error: %w", err)
	}
	refreshToken, expiresAt, err := uc.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("issue tokens generate refresh error: %w", err)
	}

	storedToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hasher.HashToken(refreshToken),
		ExpiresAt: expiresAt,
	}

	if err := uc.tokenRepo.Create(ctx, storedToken); err != nil {
		return nil, fmt.Errorf("Issue tokens save refresh: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  access_token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (uc *authUseCase) Login(ctx context.Context, input domain.LoginInput) (*domain.AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("Login get user error: %w", err)
	}

	if !hasher.CheckPassword(input.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	return uc.issueTokens(ctx, user)
}

func (uc *authUseCase) Refresh(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	if err := uc.jwtService.ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}
	tokenHash := hasher.HashToken(refreshToken)
	storedToken, err := uc.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if time.Now().After(storedToken.ExpiresAt) {
		_ = uc.tokenRepo.DeleteByHash(ctx, tokenHash)
		return nil, domain.ErrExpiredToken
	}

	user, err := uc.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("Refresh get user error: %w", err)
	}
	if err := uc.tokenRepo.DeleteByHash(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("Refresh delete old token error: %w", err)
	}
	return uc.issueTokens(ctx, user)
}

func (uc *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hasher.HashToken(refreshToken)
	if err := uc.tokenRepo.DeleteByHash(ctx, tokenHash); err != nil {
		if err == domain.ErrInvalidToken {
			return nil
		}
		return fmt.Errorf("Logout error: %w", err)
	}

	return nil
}

func (uc *authUseCase) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	if err := uc.tokenRepo.DeleteAllByUserID(ctx, userID); err != nil {
		return fmt.Errorf("logout all error: %w", err)
	}

	return nil
}

func (uc *authUseCase) VerifyEmail(ctx context.Context, token string) error {
	user, err := uc.userRepo.GetByVerifyToken(ctx, token)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("Verify email get user error: %w", err)
	}

	if user.IsVerified {
		return domain.ErrAlreadyVerified
	}

	if time.Now().After(user.VerifyTokenExp) {
		return domain.ErrExpiredToken
	}

	user.IsVerified = true
	user.VerifyToken = ""
	user.VerifyTokenExp = time.Time{}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("Verify email update user error: %w", err)
	}

	return nil
}

func (uc *authUseCase) ResendVerifiction(ctx context.Context, email string) error {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil
		}
		return fmt.Errorf("Resend verification error: %w", err)
	}

	if user.IsVerified {
		return domain.ErrUserAlreadyExsists
	}

	verifyToken, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("resend generate token: %w", err)
	}

	user.VerifyToken = verifyToken
	user.VerifyTokenExp = time.Now().Add(24 * time.Hour)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("Resend update user error: %w", err)
	}

	if err := uc.email.SendVerification(user.Email, verifyToken); err != nil {
		return fmt.Errorf("resend send email error: %w", err)
	}

	return nil

}

func (uc *authUseCase) GetProfile(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("Generate secure token error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
