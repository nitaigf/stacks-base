package services

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
)

type AuthResult struct {
	User             repositories.User
	AccessToken      string
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type AuthService struct {
	repo   repositories.Repository
	tokens *TokenService
	email  EmailSender
}

func NewAuthService(repo repositories.Repository, tokens *TokenService, email EmailSender) *AuthService {
	return &AuthService{repo: repo, tokens: tokens, email: email}
}

// REF.AUTH-01|Register
func (s *AuthService) Register(ctx context.Context, input schemas.RegisterRequest) (AuthResult, error) {
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return AuthResult{}, InternalError("failed to hash password", err)
	}

	var result AuthResult
	registrationEmail := strings.ToLower(strings.TrimSpace(input.Email))
	registrationName := strings.TrimSpace(input.Name)

	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		if _, err := store.FindUserByEmail(ctx, registrationEmail); err == nil {
			return NewAppError(http.StatusConflict, "email_taken", "email already exists", nil)
		} else if !repositories.IsNotFound(err) {
			return InternalError("failed to inspect existing user", err)
		}

		user, err := store.CreateUser(ctx, registrationName, registrationEmail, passwordHash)
		if err != nil {
			return InternalError("failed to create user", err)
		}

		issued, err := s.issueTokens(ctx, store, user)
		if err != nil {
			return err
		}

		result = issued
		return nil
	}); err != nil {
		return AuthResult{}, err
	}

	go s.audit("auth.register", &result.User.ID, "users", result.User.ID, map[string]any{"email": result.User.Email})
	go s.sendRegistrationEmail(result.User)

	return result, nil
}

// REF.AUTH-02|Login
func (s *AuthService) Login(ctx context.Context, input schemas.LoginRequest) (AuthResult, error) {
	user, err := s.repo.FindUserByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if repositories.IsNotFound(err) {
			return AuthResult{}, NewAppError(http.StatusUnauthorized, "invalid_credentials", "email or password is invalid", err)
		}

		return AuthResult{}, InternalError("failed to load user credentials", err)
	}

	valid, err := VerifyPassword(user.PasswordHash, input.Password)
	if err != nil {
		return AuthResult{}, InternalError("failed to verify password", err)
	}

	if !valid {
		return AuthResult{}, NewAppError(http.StatusUnauthorized, "invalid_credentials", "email or password is invalid", nil)
	}

	var result AuthResult
	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		issued, err := s.issueTokens(ctx, store, user)
		if err != nil {
			return err
		}

		result = issued
		return nil
	}); err != nil {
		return AuthResult{}, err
	}

	go s.audit("auth.login", &user.ID, "users", user.ID, map[string]any{"email": user.Email})

	return result, nil
}

// REF.AUTH-03|Logout
func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return NewAppError(http.StatusUnauthorized, "unauthorized", "refresh cookie is required", nil)
	}

	if _, err := s.tokens.ParseRefreshToken(rawRefreshToken); err != nil {
		return NewAppError(http.StatusUnauthorized, "unauthorized", "invalid refresh token", err)
	}

	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		if err := store.RevokeRefreshToken(ctx, HashToken(rawRefreshToken)); err != nil {
			if repositories.IsNotFound(err) {
				return NewAppError(http.StatusUnauthorized, "unauthorized", "session not found", err)
			}

			return InternalError("failed to revoke refresh token", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// REF.AUTH-04|Me
func (s *AuthService) Me(ctx context.Context, userID string) (repositories.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		if repositories.IsNotFound(err) {
			return repositories.User{}, NewAppError(http.StatusUnauthorized, "unauthorized", "failed to resolve user", err)
		}

		return repositories.User{}, InternalError("failed to load authenticated user", err)
	}

	return user, nil
}

// REF.AUTH-05|ParseAccessToken
func (s *AuthService) ParseAccessToken(raw string) (AuthClaims, error) {
	return s.tokens.ParseAccessToken(raw)
}

func (s *AuthService) issueTokens(ctx context.Context, store repositories.Repository, user repositories.User) (AuthResult, error) {
	accessToken, refreshToken, refreshExpiresAt, err := s.tokens.IssueTokens(user)
	if err != nil {
		return AuthResult{}, InternalError("failed to issue session tokens", err)
	}

	if err := store.SaveRefreshToken(ctx, user.ID, HashToken(refreshToken), refreshExpiresAt); err != nil {
		return AuthResult{}, InternalError("failed to persist refresh token", err)
	}

	return AuthResult{
		User:             user,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func (s *AuthService) audit(action string, actorUserID *string, resource string, resourceID string, metadata map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.repo.CreateAuditLog(ctx, actorUserID, action, resource, resourceID, metadata)
}

func (s *AuthService) sendRegistrationEmail(user repositories.User) {
	if s.email == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.email.SendRegistrationNotice(ctx, user); err != nil {
		log.Printf("failed to send registration email to %s: %v", user.Email, err)
	}
}