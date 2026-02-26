package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           Role      `json:"role"`
	IsVerified     bool      `json:"is_verified"`
	VerifyToken    string    `json:"-"`
	VerifyTokenExp time.Time `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegisterInput struct {
	Email    string `json:"email"   binding:"required,email"`
	Password string `json:"password" binging:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email"   binding:"required,email"`
	Password string `json:"password" binging:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}
