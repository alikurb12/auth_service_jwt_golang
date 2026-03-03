package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alikurb12/auth_service_jwt_golang/internal/handler"
	"github.com/alikurb12/auth_service_jwt_golang/internal/infrastructure/jwt"
	infraPostgres "github.com/alikurb12/auth_service_jwt_golang/internal/infrastructure/postgres"
	repoPostgres "github.com/alikurb12/auth_service_jwt_golang/internal/repository/postgres"
	"github.com/alikurb12/auth_service_jwt_golang/internal/usecase"
	"github.com/alikurb12/auth_service_jwt_golang/pkg/config"
	"github.com/alikurb12/auth_service_jwt_golang/pkg/email"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting auth-service on port %s [%s]", cfg.App.Port, cfg.App.Env)

	ctx := context.Background()
	db, err := infraPostgres.NewPool(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	userRepo := repoPostgres.NewUserRepository(db)
	tokenRepo := repoPostgres.NewTokenRepository(db)

	jwtService := jwt.NewJWTService(cfg.JWT)
	emailSender := email.NewSender(cfg.Email)

	authUC := usecase.NewAuthUseCase(
		userRepo,
		tokenRepo,
		jwtService,
		emailSender,
	)

	router := handler.NewRouter(authUC, jwtService)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://localhost:%s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}