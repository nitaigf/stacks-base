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
		if err := config.ApplyMigrationDir(ctx, db, "migrations"); err != nil {
			return fmt.Errorf("apply migrations: %w", err)
		}
	}

	if err := ensureInitialAdmin(ctx, db, cfg); err != nil {
		return fmt.Errorf("seed initial admin: %w", err)
	}
	if err := ensureDemoData(ctx, db, cfg); err != nil {
		return fmt.Errorf("seed demo data: %w", err)
	}

	repo := repositories.NewPostgresRepository(db)
	tokenService, err := services.NewTokenService(cfg)
	if err != nil {
		return fmt.Errorf("create token service: %w", err)
	}
	emailService := services.NewSMTPEmailService(cfg)
	authService := services.NewAuthService(repo, tokenService, emailService, cfg.FrontendBaseURL)
	userService := services.NewUserService(repo)
	auditService := services.NewAuditService(repo)

	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           routes.New(cfg, authService, userService, auditService),
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
	if err := db.QueryRowContext(ctx, `select exists(select 1 from users where lower(email) = lower($1) and deleted_at is null)`, email).Scan(&exists); err != nil {
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

func ensureDemoData(ctx context.Context, db *sql.DB, cfg config.Config) error {
	if !cfg.DemoSeedEnabled {
		return nil
	}

	passwordHash, err := services.HashPassword(cfg.DemoSeedPassword)
	if err != nil {
		return fmt.Errorf("hash demo seed password: %w", err)
	}

	type demoUser struct {
		Name      string
		Email     string
		Role      string
		Status    string
		DeletedAt bool
	}

	users := []demoUser{
		{Name: "Diego Owner", Email: "diego.owner@stacks-base.local", Role: "admin", Status: "active"},
		{Name: "Marina Silva", Email: "marina.silva@stacks-base.local", Role: "member", Status: "active"},
		{Name: "Bruno Costa", Email: "bruno.costa@stacks-base.local", Role: "member", Status: "active"},
		{Name: "Carla Mendes", Email: "carla.mendes@stacks-base.local", Role: "member", Status: "inactive"},
		{Name: "Helena Souza", Email: "helena.souza@stacks-base.local", Role: "member", Status: "inactive", DeletedAt: true},
		{Name: "Paulo Lima", Email: "paulo.lima@stacks-base.local", Role: "member", Status: "active"},
	}

	for _, user := range users {
		var exists bool
		if err := db.QueryRowContext(ctx, `select exists(select 1 from users where lower(email) = lower($1))`, user.Email).Scan(&exists); err != nil {
			return fmt.Errorf("check demo user %s: %w", user.Email, err)
		}
		if exists {
			continue
		}

		if _, err := db.ExecContext(ctx, `
			insert into users (name, email, password_hash, role, status, deleted_at)
			values ($1, $2, $3, $4, $5, case when $6 then now() else null end)
		`, user.Name, strings.ToLower(user.Email), passwordHash, user.Role, user.Status, user.DeletedAt); err != nil {
			return fmt.Errorf("insert demo user %s: %w", user.Email, err)
		}
	}

	var auditCount int
	if err := db.QueryRowContext(ctx, `select count(*) from audit_logs`).Scan(&auditCount); err != nil {
		return fmt.Errorf("count audit logs: %w", err)
	}
	if auditCount > 0 {
		return nil
	}

	rows, err := db.QueryContext(ctx, `
		select id, name, email
		from users
		where deleted_at is null
		order by created_at asc
		limit 3
	`)
	if err != nil {
		return fmt.Errorf("load users for audit seeds: %w", err)
	}
	defer rows.Close()

	type actor struct {
		ID    string
		Name  string
		Email string
	}

	actors := make([]actor, 0, 3)
	for rows.Next() {
		var current actor
		if err := rows.Scan(&current.ID, &current.Name, &current.Email); err != nil {
			return fmt.Errorf("scan audit seed user: %w", err)
		}
		actors = append(actors, current)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate audit seed users: %w", err)
	}
	if len(actors) == 0 {
		return nil
	}

	for index, actor := range actors {
		createdAt := time.Now().Add(time.Duration(-index-1) * time.Hour)
		if _, err := db.ExecContext(ctx, `
			insert into audit_logs (
				actor_user_id,
				actor_name,
				actor_email,
				action,
				resource,
				resource_id,
				route,
				method,
				ip_address,
				user_agent,
				metadata,
				created_at
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb, $12)
		`,
			actor.ID,
			actor.Name,
			actor.Email,
			"users.list",
			"users",
			actor.ID,
			"/api/v1/users",
			"GET",
			"127.0.0.1",
			"demo-seed",
			`{"seed":true}`,
			createdAt,
		); err != nil {
			return fmt.Errorf("insert audit seed for %s: %w", actor.Email, err)
		}
	}

	return nil
}
