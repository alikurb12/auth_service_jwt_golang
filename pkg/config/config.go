package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	JWT      JWTConfig
	Email    EmailConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type PostgresConfig struct {
	DSN string //url for connection
}

type JWTConfig struct {
	AccessSecret  string
	RefresjSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
	AppURL       string
}

func load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Postgres: PostgresConfig{
			DSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  mustGetEnv("JWT_ACCESS_SECRET"),
			RefresjSecret: mustGetEnv("JWT_REFRESH_SECRET"),
			AccessTTL:     parseDuration(getEnv("JWT_ACCESS_TTL", "15m")),
			RefreshTTL:    parseDuration(getEnv("JWT_REGRESH_TTL", "720h")),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     smtpPort,
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromAddress:  getEnv("SMTP_FROM", "gurbansahedovali@gmail.com"),
			AppURL:       getEnv("APP_URL", "http://localhost:8080"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Required env variable %v is not set", v)
	}
	return v
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Invalid duration format: %s", s)
	}
	return d
}
