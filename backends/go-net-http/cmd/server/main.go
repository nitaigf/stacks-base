package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"stacks-base/backends/go-net-http/internal/config"
	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/routes"
	"stacks-base/backends/go-net-http/internal/services"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	if cfg.AutoMigrate {
		if err := config.ApplyMigrationFile(ctx, db, "migrations/001_init.up.sql"); err != nil {
			return fmt.Errorf("apply migrations: %w", err)
		}
	}

	if err := ensureInitialAdmin(ctx, db, cfg); err != nil {
		return fmt.Errorf("seed initial admin: %w", err)
	}

	repo := repositories.NewPostgresRepository(db)
	tokenService, err := services.NewTokenService(cfg)
	if err != nil {
		return fmt.Errorf("create token service: %w", err)
	}
	emailService := services.NewSMTPEmailService(cfg)
	authService := services.NewAuthService(repo, tokenService, emailService)

	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           routes.New(cfg, authService),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownErrors := make(chan error, 1)
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownErrors <- server.Shutdown(ctx)
	}()

	log.Printf("go-net-http backend listening on %s", cfg.Address())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	if err := <-shutdownErrors; err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}

func ensureInitialAdmin(ctx context.Context, db *sql.DB, cfg config.Config) error {
	if !cfg.AdminSeedEnabled {
		return nil
	}

	email := strings.ToLower(strings.TrimSpace(cfg.AdminInitialEmail))
	name := strings.TrimSpace(cfg.AdminInitialName)
	password := cfg.AdminInitialPass

	var exists bool
	if err := db.QueryRowContext(ctx, `select exists(select 1 from users where email = $1)`, email).Scan(&exists); err != nil {
		return fmt.Errorf("check existing admin by email: %w", err)
	}

	if exists {
		log.Printf("initial admin already exists for email %s", email)
		return nil
	}

	passwordHash, err := services.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash initial admin password: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
		insert into users (name, email, password_hash, role, status)
		values ($1, $2, $3, 'admin', 'active')
	`, name, email, passwordHash); err != nil {
		return fmt.Errorf("insert initial admin user: %w", err)
	}

	log.Printf("initial admin seeded: %s", email)
	return nil
}