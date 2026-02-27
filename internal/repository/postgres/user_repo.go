package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
			INSERT INTO users (id, email, password_hash, role, is_verified, verify_token, verify_token_exp)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		u.ID,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.IsVerified,
		u.VerifyToken,
		u.VerifyTokenExp,
	).Scan(&u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrUserAlreadyExsists
		}
		return fmt.Errorf("Create user error: %v", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_verified, verify_token,
			   verify_token_exp, created_at, updated_at
		FROM user
		WHERE id=$1`

	return r.scanUser(r.db.QueryRow(ctx, query, id))
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_verified, verify_token,
			   verify_token_exp, created_at, updated_at
		FROM user
		WHERE email=$1`

	return r.scanUser(r.db.QueryRow(ctx, query, email))
}

func (r *userRepository) GetByVerifyToken(ctx context.Context, token string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_verified, verify_token,
			   verify_token_exp, created_at, updated_at
		FROM user
		WHERE verify_token=$1`

	return r.scanUser(r.db.QueryRow(ctx, query, token))
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users
		SET email            = $1,
		    password_hash    = $2,
		    role             = $3,
		    is_verified      = $4,
		    verify_token     = $5,
		    verify_token_exp = $6,
		    updated_at       = NOW()
		WHERE id = $7
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.IsVerified,
		u.VerifyToken,
		u.VerifyTokenExp,
		u.ID,
	).Scan(&u.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("Update user error: %v", err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id=$1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Delete user error: %v", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) scanUser(row pgx.Row) (*domain.User, error) {
	u := &domain.User{}

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.IsVerified,
		&u.VerifyToken,
		&u.VerifyTokenExp,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("Scan user error: %v", err)
	}

	return u, nil
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
