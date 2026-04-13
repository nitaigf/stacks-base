package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
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
	nextUserID    int

	createUserErr              error
	updateUserErr              error
	setUserStatusErr           error
	softDeleteUserErr          error
	restoreUserErr             error
	hardDeleteUserErr          error
	updateUserPasswordErr      error
	findUserByEmailErr         error
	findUserByIDErr            error
	listUsersErr               error
	recordUserLoginErr         error
	saveRefreshTokenErr        error
	revokeRefreshTokenErr      error
	revokeRefreshTokensByIDErr error
	savePasswordResetTokenErr  error
	findPasswordResetTokenErr  error
	usePasswordResetTokenErr   error
	revokeActiveResetTokensErr error
	createAuditLogErr          error
	listAuditLogsErr           error
	runInTxErr                 error
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
		nextUserID:    1,
	}
}

func (r *memoryRepository) CreateUser(_ context.Context, params repositories.CreateUserParams) (repositories.User, error) {
	if r.createUserErr != nil {
		return repositories.User{}, r.createUserErr
	}

	user := repositories.User{
		ID:           r.nextID(),
		Name:         strings.TrimSpace(params.Name),
		Email:        normalizeTestEmail(params.Email),
		PasswordHash: params.PasswordHash,
		Role:         params.Role,
		Status:       params.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return r.storeUser(user), nil
}

func (r *memoryRepository) UpdateUser(_ context.Context, id string, params repositories.UpdateUserParams) (repositories.User, error) {
	if r.updateUserErr != nil {
		return repositories.User{}, r.updateUserErr
	}

	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	previousEmail := normalizeTestEmail(user.Email)
	user.Name = strings.TrimSpace(params.Name)
	user.Email = normalizeTestEmail(params.Email)
	user.Role = params.Role
	user.UpdatedAt = time.Now()
	if previousEmail != user.Email {
		delete(r.users, previousEmail)
	}
	return r.storeUser(user), nil
}

func (r *memoryRepository) SetUserStatus(_ context.Context, id string, status string) (repositories.User, error) {
	if r.setUserStatusErr != nil {
		return repositories.User{}, r.setUserStatusErr
	}

	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.Status = status
	user.UpdatedAt = time.Now()
	return r.storeUser(user), nil
}

func (r *memoryRepository) SoftDeleteUser(_ context.Context, id string, deletedBy *string) (repositories.User, error) {
	if r.softDeleteUserErr != nil {
		return repositories.User{}, r.softDeleteUserErr
	}

	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	now := time.Now()
	user.DeletedAt = &now
	user.DeletedBy = deletedBy
	user.Status = "inactive"
	user.UpdatedAt = now
	return r.storeUser(user), nil
}

func (r *memoryRepository) RestoreUser(_ context.Context, id string) (repositories.User, error) {
	if r.restoreUserErr != nil {
		return repositories.User{}, r.restoreUserErr
	}

	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.DeletedAt = nil
	user.DeletedBy = nil
	user.Status = "active"
	user.UpdatedAt = time.Now()
	return r.storeUser(user), nil
}

func (r *memoryRepository) HardDeleteUser(_ context.Context, id string) error {
	if r.hardDeleteUserErr != nil {
		return r.hardDeleteUserErr
	}

	user, ok := r.byID[id]
	if !ok {
		return sql.ErrNoRows
	}
	delete(r.byID, id)
	delete(r.users, normalizeTestEmail(user.Email))
	return nil
}

func (r *memoryRepository) UpdateUserPassword(_ context.Context, id string, passwordHash string) (repositories.User, error) {
	if r.updateUserPasswordErr != nil {
		return repositories.User{}, r.updateUserPasswordErr
	}

	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now()
	return r.storeUser(user), nil
}

func (r *memoryRepository) FindUserByEmail(_ context.Context, email string, includeDeleted bool) (repositories.User, error) {
	if r.findUserByEmailErr != nil {
		return repositories.User{}, r.findUserByEmailErr
	}

	user, ok := r.users[normalizeTestEmail(email)]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	if !includeDeleted && user.DeletedAt != nil {
		return repositories.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (r *memoryRepository) FindUserByID(_ context.Context, id string, includeDeleted bool) (repositories.User, error) {
	if r.findUserByIDErr != nil {
		return repositories.User{}, r.findUserByIDErr
	}

	user, ok := r.byID[id]
	if !ok {
		return repositories.User{}, sql.ErrNoRows
	}
	if !includeDeleted && user.DeletedAt != nil {
		return repositories.User{}, sql.ErrNoRows
	}
	return user, nil
}

func (r *memoryRepository) ListUsers(_ context.Context, params repositories.UserListParams) ([]repositories.User, repositories.PaginationMeta, error) {
	if r.listUsersErr != nil {
		return nil, repositories.PaginationMeta{}, r.listUsersErr
	}

	users := make([]repositories.User, 0, len(r.byID))
	for _, user := range r.byID {
		if !params.IncludeDeleted && user.DeletedAt != nil {
			continue
		}
		users = append(users, user)
	}
	return users, repositories.PaginationMeta{Page: 1, PerPage: len(users), Total: len(users), TotalPages: 1}, nil
}

func (r *memoryRepository) RecordUserLogin(_ context.Context, id string, at time.Time) error {
	if r.recordUserLoginErr != nil {
		return r.recordUserLoginErr
	}

	user, ok := r.byID[id]
	if !ok {
		return sql.ErrNoRows
	}
	user.LastLoginAt = &at
	r.storeUser(user)
	return nil
}

func (r *memoryRepository) SaveRefreshToken(_ context.Context, userID string, tokenHash string, _ time.Time) error {
	if r.saveRefreshTokenErr != nil {
		return r.saveRefreshTokenErr
	}

	r.refreshTokens[tokenHash] = userID
	return nil
}

func (r *memoryRepository) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	if r.revokeRefreshTokenErr != nil {
		return r.revokeRefreshTokenErr
	}

	if _, ok := r.refreshTokens[tokenHash]; !ok {
		return sql.ErrNoRows
	}
	delete(r.refreshTokens, tokenHash)
	return nil
}

func (r *memoryRepository) RevokeRefreshTokensByUserID(_ context.Context, userID string) error {
	if r.revokeRefreshTokensByIDErr != nil {
		return r.revokeRefreshTokensByIDErr
	}

	for tokenHash, currentUserID := range r.refreshTokens {
		if currentUserID == userID {
			delete(r.refreshTokens, tokenHash)
		}
	}
	return nil
}

func (r *memoryRepository) SavePasswordResetToken(_ context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	if r.savePasswordResetTokenErr != nil {
		return r.savePasswordResetTokenErr
	}

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
	if r.findPasswordResetTokenErr != nil {
		return repositories.PasswordResetToken{}, r.findPasswordResetTokenErr
	}

	token, ok := r.resetTokens[tokenHash]
	if !ok {
		return repositories.PasswordResetToken{}, sql.ErrNoRows
	}
	return token, nil
}

func (r *memoryRepository) UsePasswordResetToken(_ context.Context, tokenHash string) error {
	if r.usePasswordResetTokenErr != nil {
		return r.usePasswordResetTokenErr
	}

	token, ok := r.resetTokens[tokenHash]
	if !ok {
		return sql.ErrNoRows
	}
	now := time.Now()
	token.UsedAt = &now
	r.resetTokens[tokenHash] = token
	return nil
}

func (r *memoryRepository) RevokeActivePasswordResetTokens(_ context.Context, userID string) error {
	if r.revokeActiveResetTokensErr != nil {
		return r.revokeActiveResetTokensErr
	}
	now := time.Now()
	for hash, token := range r.resetTokens {
		if token.UserID == userID && token.UsedAt == nil && token.ExpiresAt.After(now) {
			token.UsedAt = &now
			r.resetTokens[hash] = token
		}
	}
	return nil
}

func (r *memoryRepository) CreateAuditLog(_ context.Context, _ repositories.AuditLogInput) error {
	return r.createAuditLogErr
}

func (r *memoryRepository) ListAuditLogs(_ context.Context, _ repositories.AuditLogListParams) ([]repositories.AuditLog, repositories.PaginationMeta, error) {
	return nil, repositories.PaginationMeta{}, r.listAuditLogsErr
}

func (r *memoryRepository) RunInTx(_ context.Context, fn func(repositories.Repository) error) error {
	if r.runInTxErr != nil {
		return r.runInTxErr
	}
	return fn(r)
}

func (r *memoryRepository) seedUser(user repositories.User) repositories.User {
	if user.ID == "" {
		user.ID = r.nextID()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = user.CreatedAt
	}
	if user.Status == "" {
		user.Status = "active"
	}
	if user.Role == "" {
		user.Role = "member"
	}
	user.Email = normalizeTestEmail(user.Email)
	return r.storeUser(user)
}

func (r *memoryRepository) storeUser(user repositories.User) repositories.User {
	user.Email = normalizeTestEmail(user.Email)
	r.byID[user.ID] = user
	r.users[user.Email] = user
	return user
}

func (r *memoryRepository) nextID() string {
	id := fmt.Sprintf("user-%d", r.nextUserID)
	r.nextUserID++
	return id
}

func normalizeTestEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
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

func TestAuthServiceRegisterRejectsExistingEmail(t *testing.T) {
	repo := newMemoryRepository()
	tokens, _ := NewTokenService(config.Config{AccessTokenTTL: time.Hour, RefreshTokenTTL: time.Hour})
	service := NewAuthService(repo, tokens, noopEmailSender{}, "http://localhost:3000")

	repo.seedUser(repositories.User{Email: "nitai@example.com"})

	_, err := service.Register(context.Background(), RequestMetadata{}, schemas.RegisterRequest{
		Name:     "Nitai",
		Email:    "NITAI@example.com", // testing case insensitivity
		Password: "password123",
	})
	assertAppErrorCode(t, err, http.StatusConflict, "email_taken")
}

func TestAuthServiceLoginRejectsInactiveUser(t *testing.T) {
	repo := newMemoryRepository()
	tokens, _ := NewTokenService(config.Config{AccessTokenTTL: time.Hour, RefreshTokenTTL: time.Hour})
	service := NewAuthService(repo, tokens, noopEmailSender{}, "http://localhost:3000")

	passwordHash, _ := HashPassword("password123")
	repo.seedUser(repositories.User{Email: "nitai@example.com", PasswordHash: passwordHash, Status: "inactive"})

	_, err := service.Login(context.Background(), RequestMetadata{}, schemas.LoginRequest{
		Email:    "nitai@example.com",
		Password: "password123",
	})
	assertAppErrorCode(t, err, http.StatusForbidden, "account_inactive")
}

func TestAuthServiceForgotPasswordIgnoresNonexistentEmail(t *testing.T) {
	repo := newMemoryRepository()
	service := NewAuthService(repo, nil, noopEmailSender{}, "http://localhost:3000")

	err := service.ForgotPassword(context.Background(), RequestMetadata{}, schemas.ForgotPasswordRequest{
		Email: "nonexistent@example.com",
	})
	if err != nil {
		t.Fatalf("expected success (silent), got %v", err)
	}
	if len(repo.resetTokens) > 0 {
		t.Fatalf("expected no reset token to be created")
	}
}

func TestAuthServiceForgotPasswordInvalidatesOldTokens(t *testing.T) {
	repo := newMemoryRepository()
	service := NewAuthService(repo, nil, noopEmailSender{}, "http://localhost:3000")

	user := repo.seedUser(repositories.User{ID: "user-1", Email: "nitai@example.com", Status: "active"})
	
	// Issue first token
	_ = service.ForgotPassword(context.Background(), RequestMetadata{}, schemas.ForgotPasswordRequest{Email: user.Email})
	if len(repo.resetTokens) != 1 {
		t.Fatalf("expected 1 reset token")
	}

	// Issue second token
	_ = service.ForgotPassword(context.Background(), RequestMetadata{}, schemas.ForgotPasswordRequest{Email: user.Email})
	if len(repo.resetTokens) != 2 {
		t.Fatalf("expected 2 reset tokens in total")
	}

	// Check that the first one is used (invalidated)
	usedCount := 0
	for _, t := range repo.resetTokens {
		if t.UsedAt != nil {
			usedCount++
		}
	}
	if usedCount != 1 {
		t.Fatalf("expected 1 token to be invalidated, got %d", usedCount)
	}
}
