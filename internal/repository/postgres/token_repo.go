package postgres

import (
	"context"
	"fmt"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *tokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, t *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	err := r.db.QueryRow(ctx, query,
		t.ID,
		t.UserID,
		t.TokenHash,
		t.ExpiresAt,
	).Scan(&t.CreatedAt)

	if err != nil {
		return fmt.Errorf("Create refresh token error: %v", err)
	}

	return nil
}

func (r *tokenRepository) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at FROM refresh_tokens
			  WHERE token_hash=$1`

	t := &domain.RefreshToken{}
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.ExpiresAt,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Getting token by id error: %v", err)
	}
	return t, nil
}

func (r *tokenRepository) DeleteByHash(ctx context.Context, hash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash=$1`

	result, err := r.db.Exec(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("Deleting token by hash error: %v", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInvalidToken
	}
	return nil
}

func (r *tokenRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id=$1`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("Deleting token by user id error: %v", err)
	}

	return nil
}
