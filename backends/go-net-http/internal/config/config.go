package config

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName            string
	AppEnv             string
	BackendHost        string
	BackendPort        string
	AllowedOrigin      string
	DatabaseHost       string
	DatabasePort       string
	DatabaseName       string
	DatabaseUser       string
	DatabasePassword   string
	DatabaseSSLMode    string
	FrontendBaseURL    string
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	SMTPHost           string
	SMTPPort           string
	SMTPUsername       string
	SMTPPassword       string
	SMTPFromName       string
	SMTPFromAddress    string
	AutoMigrate        bool
	AdminSeedEnabled   bool
	AdminInitialName   string
	AdminInitialEmail  string
	AdminInitialPass   string
}

func Load() (Config, error) {
	loadEnvFiles()

	accessTTL, err := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("parse ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := time.ParseDuration(getEnv("REFRESH_TOKEN_TTL", "168h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse REFRESH_TOKEN_TTL: %w", err)
	}

	autoMigrate := true
	if raw := strings.TrimSpace(os.Getenv("AUTO_MIGRATE")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse AUTO_MIGRATE: %w", err)
		}
		autoMigrate = parsed
	}

	adminSeedEnabled := true
	if raw := strings.TrimSpace(os.Getenv("ADMIN_SEED_ENABLED")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse ADMIN_SEED_ENABLED: %w", err)
		}
		adminSeedEnabled = parsed
	}

	cfg := Config{
		AppName:            getEnv("APP_NAME", "stacks-base"),
		AppEnv:             getEnv("APP_ENV", "development"),
		BackendHost:        getEnv("BACKEND_HOST", "127.0.0.1"),
		BackendPort:        getEnv("BACKEND_PORT", "8080"),
		AllowedOrigin:      getEnv("BACKEND_ALLOWED_ORIGIN", "http://127.0.0.1:3000"),
		DatabaseHost:       getEnv("DATABASE_HOST", "127.0.0.1"),
		DatabasePort:       getEnv("DATABASE_PORT", "5432"),
		DatabaseName:       getEnv("DATABASE_NAME", "stacks_base"),
		DatabaseUser:       getEnv("DATABASE_USER", "postgres"),
		DatabasePassword:   getOptionalEnv("DATABASE_PASSWORD", "postgres"),
		DatabaseSSLMode:    getEnv("DATABASE_SSLMODE", "disable"),
		FrontendBaseURL:    getEnv("FRONTEND_BASE_URL", "http://127.0.0.1:3000"),
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "replace-me-with-a-long-random-secret"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "replace-me-with-a-second-long-random-secret"),
		AccessTokenTTL:     accessTTL,
		RefreshTokenTTL:    refreshTTL,
		SMTPHost:           getEnv("SMTP_HOST", ""),
		SMTPPort:           getEnv("SMTP_PORT", "2525"),
		SMTPUsername:       getEnv("SMTP_USERNAME", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		SMTPFromName:       getEnv("SMTP_FROM_NAME", "Stacks Base"),
		SMTPFromAddress:    getEnv("SMTP_FROM_ADDRESS", "no-reply@stacks-base.local"),
		AutoMigrate:        autoMigrate,
		AdminSeedEnabled:   adminSeedEnabled,
		AdminInitialName:   getEnv("ADMIN_INITIAL_NAME", "Admin"),
		AdminInitialEmail:  getEnv("ADMIN_INITIAL_EMAIL", "admin@stacks-base.local"),
		AdminInitialPass:   getEnv("ADMIN_INITIAL_PASSWORD", "Admin@123456"),
	}

	if cfg.AccessTokenSecret == "" || cfg.RefreshTokenSecret == "" {
		return Config{}, fmt.Errorf("token secrets are required")
	}

	if cfg.SMTPHost == "" || cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
		return Config{}, fmt.Errorf("smtp configuration is required")
	}

	if cfg.AdminSeedEnabled {
		if strings.TrimSpace(cfg.AdminInitialEmail) == "" {
			return Config{}, fmt.Errorf("ADMIN_INITIAL_EMAIL is required when ADMIN_SEED_ENABLED=true")
		}

		if strings.TrimSpace(cfg.AdminInitialPass) == "" {
			return Config{}, fmt.Errorf("ADMIN_INITIAL_PASSWORD is required when ADMIN_SEED_ENABLED=true")
		}
	}

	return cfg, nil
}

func (c Config) Address() string {
	return c.BackendHost + ":" + c.BackendPort
}

func (c Config) DatabaseDSN() string {
	parts := []string{
		fmt.Sprintf("host=%s", c.DatabaseHost),
		fmt.Sprintf("port=%s", c.DatabasePort),
		fmt.Sprintf("dbname=%s", c.DatabaseName),
		fmt.Sprintf("user=%s", c.DatabaseUser),
		fmt.Sprintf("sslmode=%s", c.DatabaseSSLMode),
	}

	if c.DatabasePassword != "" {
		parts = append(parts, fmt.Sprintf("password=%s", c.DatabasePassword))
	}

	return strings.Join(parts, " ")
}

func ApplyMigrationFile(ctx context.Context, db *sql.DB, relativePath string) error {
	base, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	contents, err := os.ReadFile(filepath.Join(base, relativePath))
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	if _, err := db.ExecContext(ctx, string(contents)); err != nil {
		return fmt.Errorf("execute migration file: %w", err)
	}

	return nil
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return value
	}

	return fallback
}

func getOptionalEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func loadEnvFiles() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return
	}

	for _, candidate := range []string{
		filepath.Join(workingDirectory, ".env"),
		filepath.Join(workingDirectory, "..", ".env"),
		filepath.Join(workingDirectory, "..", "..", ".env"),
	} {
		_ = loadEnvFile(candidate)
	}
}

func loadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		value = strings.Trim(strings.TrimSpace(value), `"`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}