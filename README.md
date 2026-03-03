<div align="center">

# 🔐 Auth Service

**Production-ready JWT Authentication & Authorization REST API**

Built with Go · Gin · PostgreSQL · Clean Architecture

[![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-1.10-008ECF?style=flat)](https://github.com/gin-gonic/gin)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-golang--jwt%2Fv5-F7B93E?style=flat)](https://github.com/golang-jwt/jwt)

</div>

---

## ✨ Features

- **Registration & Login** — secure user registration and authentication
- **JWT Access Tokens** — short-lived (15m) signed with HS256
- **Refresh Token Rotation** — long-lived refresh tokens stored as SHA-256 hashes in DB; rotated on every use
- **Email Verification** — token-based email confirmation with 24h expiry
- **RBAC** — Role-Based Access Control (`user` / `admin`)
- **Logout & Logout All** — invalidate one session or all devices at once
- **Graceful Shutdown** — clean server shutdown on SIGINT/SIGTERM
- **Clean Architecture** — Domain → Repository → UseCase → Handler layers with no leaking dependencies

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│                  HTTP Layer                  │
│         Gin Router · Handlers · Middleware   │
└────────────────────┬────────────────────────┘
                     │
┌────────────────────▼────────────────────────┐
│                 UseCase Layer                │
│          Business Logic · Auth Flow          │
└────────────────────┬────────────────────────┘
                     │
┌────────────────────▼────────────────────────┐
│              Repository Layer                │
│          Interfaces · pgx Postgres           │
└────────────────────┬────────────────────────┘
                     │
┌────────────────────▼────────────────────────┐
│               Infrastructure                 │
│          PostgreSQL Pool · JWT Service       │
└─────────────────────────────────────────────┘
```

Each layer depends only on the layer below it via **interfaces** — never on concrete implementations. This makes every layer independently testable.

---

## 📁 Project Structure

```
auth-service/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point, Dependency Injection
├── internal/
│   ├── domain/                      # Entities, interfaces, domain errors
│   │   ├── user.go
│   │   ├── token.go
│   │   └── errors.go
│   ├── repository/                  # Storage interfaces + pgx implementations
│   │   ├── repository.go
│   │   └── postgres/
│   │       ├── user_repo.go
│   │       └── token_repo.go
│   ├── usecase/                     # Business logic
│   │   ├── usecase.go
│   │   └── auth_usecase.go
│   ├── handler/                     # HTTP handlers, middleware, router
│   │   ├── auth_handler.go
│   │   ├── middleware.go
│   │   ├── response.go
│   │   └── router.go
│   └── infrastructure/
│       ├── postgres/postgres.go     # pgxpool connection
│       └── jwt/jwt.go               # JWT sign & validate
├── pkg/
│   ├── config/config.go             # Env-based config
│   ├── hasher/hasher.go             # bcrypt + SHA-256
│   └── email/email.go               # SMTP sender
├── Makefile
└── .env.example
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 13+

### 1. Clone & Install

```bash
git clone https://github.com/yourname/auth-service.git
cd auth-service
go mod tidy
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your values:

```dotenv
APP_PORT=8080
APP_ENV=development
APP_URL=http://localhost:8080

POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable

# Generate with: openssl rand -hex 32
JWT_ACCESS_SECRET=your-access-secret-here
JWT_REFRESH_SECRET=your-refresh-secret-here
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your@gmail.com


### Run the Server

```bash
go run ./cmd/api/main.go
```

```
Starting auth-service on port 8080 [development]
Connected to PostgreSQL
Server listening on http://localhost:8080
```

---

## 📡 API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/health` | ❌ | Server health check |

---

### Auth Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/auth/register` | ❌ | Register a new user |
| `POST` | `/auth/login` | ❌ | Login and receive tokens |
| `POST` | `/auth/refresh` | ❌ | Refresh access token |
| `GET` | `/auth/verify-email?token=` | ❌ | Verify email address |
| `POST` | `/auth/resend-verification` | ❌ | Resend verification email |
| `GET` | `/auth/me` | ✅ | Get current user profile |
| `POST` | `/auth/logout` | ✅ | Logout current session |
| `POST` | `/auth/logout-all` | ✅ | Logout from all devices |

---

### Admin Endpoints

| Method | Endpoint | Auth | Role |
|--------|----------|------|------|
| `GET` | `/admin/ping` | ✅ | `admin` |

---

### Request & Response Examples

#### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "role": "user",
      "is_verified": false,
      "created_at": "2026-03-03T12:00:00Z",
      "updated_at": "2026-03-03T12:00:00Z"
    }
  }
}
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

#### Get Profile
```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

#### Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<refresh_token>"}'
```

#### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<refresh_token>"}'
```

---

## 🔒 Security Design

| Concern | Solution |
|---------|----------|
| Password storage | `bcrypt` with cost factor 12 |
| Refresh token storage | SHA-256 hash stored in DB, raw token only sent to client |
| Token invalidation | Refresh tokens stored in DB — can be revoked any time |
| Token rotation | Old refresh token deleted on every `/refresh` call |
| Email enumeration | Login always returns `invalid credentials`, never `user not found` |
| Resend enumeration | Always returns success regardless of whether email exists |
| Secrets | Two independent JWT secrets for access and refresh tokens |

---

## ⚙️ Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_PORT` | ❌ | `8080` | HTTP server port |
| `APP_ENV` | ❌ | `development` | Environment name |
| `APP_URL` | ❌ | `http://localhost:8080` | Base URL (used in email links) |
| `POSTGRES_DSN` | ❌ | `postgres://...` | PostgreSQL connection string |
| `JWT_ACCESS_SECRET` | ✅ | — | Access token signing secret |
| `JWT_REFRESH_SECRET` | ✅ | — | Refresh token signing secret |
| `JWT_ACCESS_TTL` | ❌ | `15m` | Access token lifetime |
| `JWT_REFRESH_TTL` | ❌ | `720h` | Refresh token lifetime |
| `SMTP_HOST` | ❌ | `smtp.gmail.com` | SMTP server host |
| `SMTP_PORT` | ❌ | `587` | SMTP server port |
| `SMTP_USER` | ❌ | — | SMTP username |
| `SMTP_PASSWORD` | ❌ | — | SMTP password (empty = print to stdout) |
| `SMTP_FROM` | ❌ | — | Sender email address |


## 🧱 Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22 |
| HTTP Framework | Gin v1.10 |
| Database | PostgreSQL 16 |
| DB Driver | pgx v5 (connection pool) |
| Authentication | golang-jwt/jwt v5 |
| Password Hashing | bcrypt (cost 12) |
| UUID | google/uuid |
| Config | godotenv |

---
