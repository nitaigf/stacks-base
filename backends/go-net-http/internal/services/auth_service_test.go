package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	paseto "aidanwoods.dev/go-paseto"

	"stacks-base/backends/go-net-http/internal/config"
	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
)

type memoryRepository struct {
	users         map[string]repositories.User
	byID          map[string]repositories.User
	refreshTokens map[string]string
}

type noopEmailSender struct{}

func (noopEmailSender) SendRegistrationNotice(_ context.Context, _ repositories.User) error {
	return nil
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		users:         map[string]repositories.User{},
		byID:          map[string]repositories.User{},
		refreshTokens: map[string]string{},
	}
}

func (r *memoryRepository) CreateUser(_ context.Context, name string, email string, passwordHash string) (repositories.User, error) {
	user := repositories.User{
		ID:           "user-1",
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "member",
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	r.users[email] = user
	r.byID[user.ID] = user
	return user, nil
}

func (r *memoryRepository) FindUserByEmail(_ context.Context, email string) (repositories.User, error) {
	user, ok := r.users[email]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (r *memoryRepository) FindUserByID(_ context.Context, id string) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (r *memoryRepository) SaveRefreshToken(_ context.Context, userID string, tokenHash string, _ time.Time) error {
	r.refreshTokens[tokenHash] = userID
	return nil
}

func (r *memoryRepository) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	if _, ok := r.refreshTokens[tokenHash]; !ok {
		return sql.ErrNoRows
	}
	delete(r.refreshTokens, tokenHash)
	return nil
}

func (r *memoryRepository) CreateAuditLog(_ context.Context, _ *string, _ string, _ string, _ string, _ map[string]any) error {
	return nil
}

func (r *memoryRepository) RunInTx(_ context.Context, fn func(repositories.Repository) error) error {
	return fn(r)
}

func TestAuthServiceRegisterLoginAndLogout(t *testing.T) {
	repo := newMemoryRepository()
	tokens, err := NewTokenService(config.Config{
		AccessTokenSecret:  "access-secret",
		RefreshTokenSecret: "refresh-secret",
		AccessTokenTTL:     15 * time.Minute,
		RefreshTokenTTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("new token service: %v", err)
	}

	service := NewAuthService(repo, tokens, noopEmailSender{})

	registered, err := service.Register(context.Background(), schemas.RegisterRequest{
		Name:     "Nitai",
		Email:    "nitai@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if registered.AccessToken == "" || registered.RefreshToken == "" {
		t.Fatalf("expected issued tokens")
	}

	if _, err := tokens.ParseAccessToken(registered.AccessToken); err != nil {
		t.Fatalf("parse access token: %v", err)
	}

	loginResult, err := service.Login(context.Background(), schemas.LoginRequest{
		Email:    "nitai@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if err := service.Logout(context.Background(), loginResult.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}
}

func TestAuthServiceRejectsInvalidCredentials(t *testing.T) {
	repo := newMemoryRepository()
	tokens, err := NewTokenService(config.Config{
		AccessTokenSecret:  "access-secret",
		RefreshTokenSecret: "refresh-secret",
		AccessTokenTTL:     15 * time.Minute,
		RefreshTokenTTL:    time.Hour,
	})
	if err != nil {
		t.Fatalf("new token service: %v", err)
	}

	service := NewAuthService(repo, tokens, noopEmailSender{})
	_, err = service.Login(context.Background(), schemas.LoginRequest{
		Email:    "missing@example.com",
		Password: "password123",
	})
	appErr, ok := AsAppError(err)
	if !ok || appErr.Code != "invalid_credentials" {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestHashTokenStable(t *testing.T) {
	left := HashToken("refresh-token")
	right := HashToken("refresh-token")
	if left != right {
		t.Fatalf("expected stable token hash")
	}
}

func TestPasetoLibraryIsAvailable(t *testing.T) {
	key, err := paseto.V4SymmetricKeyFromBytes(hashKey("access-secret"))
	if err != nil {
		t.Fatalf("build key: %v", err)
	}

	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(time.Minute))
	token.SetString("userId", "user-1")

	raw := token.V4Encrypt(key, nil)
	parser := paseto.NewParserForValidNow()
	if _, err := parser.ParseV4Local(key, raw, nil); err != nil {
		t.Fatalf("parse v4 local token: %v", err)
	}
}