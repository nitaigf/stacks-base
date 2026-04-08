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
	resetTokens   map[string]repositories.PasswordResetToken
}

type noopEmailSender struct{}

func (noopEmailSender) SendRegistrationNotice(_ context.Context, _ repositories.User) error {
	return nil
}

func (noopEmailSender) SendPasswordResetNotice(_ context.Context, _ repositories.User, _ string) error {
	return nil
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		users:         map[string]repositories.User{},
		byID:          map[string]repositories.User{},
		refreshTokens: map[string]string{},
		resetTokens:   map[string]repositories.PasswordResetToken{},
	}
}

func (r *memoryRepository) CreateUser(_ context.Context, params repositories.CreateUserParams) (repositories.User, error) {
	user := repositories.User{
		ID:           "user-1",
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Role:         params.Role,
		Status:       params.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	r.users[params.Email] = user
	r.byID[user.ID] = user
	return user, nil
}

func (r *memoryRepository) UpdateUser(_ context.Context, id string, params repositories.UpdateUserParams) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.Name = params.Name
	user.Email = params.Email
	user.Role = params.Role
	user.UpdatedAt = time.Now()
	r.byID[id] = user
	r.users[user.Email] = user
	return user, nil
}

func (r *memoryRepository) SetUserStatus(_ context.Context, id string, status string) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.Status = status
	user.UpdatedAt = time.Now()
	r.byID[id] = user
	r.users[user.Email] = user
	return user, nil
}

func (r *memoryRepository) SoftDeleteUser(_ context.Context, id string, deletedBy *string) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	now := time.Now()
	user.DeletedAt = &now
	user.DeletedBy = deletedBy
	user.Status = "inactive"
	user.UpdatedAt = now
	r.byID[id] = user
	r.users[user.Email] = user
	return user, nil
}

func (r *memoryRepository) RestoreUser(_ context.Context, id string) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.DeletedAt = nil
	user.DeletedBy = nil
	user.UpdatedAt = time.Now()
	r.byID[id] = user
	r.users[user.Email] = user
	return user, nil
}

func (r *memoryRepository) HardDeleteUser(_ context.Context, id string) error {
	user, ok := r.byID[id]
	if !ok {
		return sql.ErrNoRows
	}
	delete(r.byID, id)
	delete(r.users, user.Email)
	return nil
}

func (r *memoryRepository) UpdateUserPassword(_ context.Context, id string, passwordHash string) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now()
	r.byID[id] = user
	r.users[user.Email] = user
	return user, nil
}

func (r *memoryRepository) FindUserByEmail(_ context.Context, email string, _ bool) (repositories.User, error) {
	user, ok := r.users[email]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (r *memoryRepository) FindUserByID(_ context.Context, id string, _ bool) (repositories.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (r *memoryRepository) ListUsers(_ context.Context, _ repositories.UserListParams) ([]repositories.User, repositories.PaginationMeta, error) {
	users := make([]repositories.User, 0, len(r.byID))
	for _, user := range r.byID {
		users = append(users, user)
	}
	return users, repositories.PaginationMeta{Page: 1, PerPage: len(users), Total: len(users), TotalPages: 1}, nil
}

func (r *memoryRepository) RecordUserLogin(_ context.Context, id string, at time.Time) error {
	user, ok := r.byID[id]
	if !ok {
		return sql.ErrNoRows
	}
	user.LastLoginAt = &at
	r.byID[id] = user
	r.users[user.Email] = user
	return nil
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

func (r *memoryRepository) RevokeRefreshTokensByUserID(_ context.Context, userID string) error {
	for tokenHash, currentUserID := range r.refreshTokens {
		if currentUserID == userID {
			delete(r.refreshTokens, tokenHash)
		}
	}
	return nil
}

func (r *memoryRepository) SavePasswordResetToken(_ context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	r.resetTokens[tokenHash] = repositories.PasswordResetToken{
		ID:        "reset-1",
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return nil
}

func (r *memoryRepository) FindPasswordResetToken(_ context.Context, tokenHash string) (repositories.PasswordResetToken, error) {
	token, ok := r.resetTokens[tokenHash]
	if !ok {
		return repositories.PasswordResetToken{}, sql.ErrNoRows
	}
	return token, nil
}

func (r *memoryRepository) UsePasswordResetToken(_ context.Context, tokenHash string) error {
	token, ok := r.resetTokens[tokenHash]
	if !ok {
		return sql.ErrNoRows
	}
	now := time.Now()
	token.UsedAt = &now
	r.resetTokens[tokenHash] = token
	return nil
}

func (r *memoryRepository) CreateAuditLog(_ context.Context, _ repositories.AuditLogInput) error {
	return nil
}

func (r *memoryRepository) ListAuditLogs(_ context.Context, _ repositories.AuditLogListParams) ([]repositories.AuditLog, repositories.PaginationMeta, error) {
	return nil, repositories.PaginationMeta{}, nil
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

	service := NewAuthService(repo, tokens, noopEmailSender{}, "http://127.0.0.1:3000")

	registered, err := service.Register(context.Background(), RequestMetadata{}, schemas.RegisterRequest{
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

	loginResult, err := service.Login(context.Background(), RequestMetadata{}, schemas.LoginRequest{
		Email:    "nitai@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if err := service.Logout(context.Background(), RequestMetadata{ActorUserID: &loginResult.User.ID}, loginResult.RefreshToken); err != nil {
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

	service := NewAuthService(repo, tokens, noopEmailSender{}, "http://127.0.0.1:3000")
	_, err = service.Login(context.Background(), RequestMetadata{}, schemas.LoginRequest{
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
