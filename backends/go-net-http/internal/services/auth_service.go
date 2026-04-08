package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	repo            repositories.Repository
	tokens          *TokenService
	email           EmailSender
	frontendBaseURL string
}

func NewAuthService(repo repositories.Repository, tokens *TokenService, email EmailSender, frontendBaseURL string) *AuthService {
	return &AuthService{
		repo:            repo,
		tokens:          tokens,
		email:           email,
		frontendBaseURL: strings.TrimRight(strings.TrimSpace(frontendBaseURL), "/"),
	}
}

// REF.AUTH-01|Register
func (s *AuthService) Register(ctx context.Context, request RequestMetadata, input schemas.RegisterRequest) (AuthResult, error) {
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return AuthResult{}, InternalError("failed to hash password", err)
	}

	var result AuthResult
	registrationEmail := strings.ToLower(strings.TrimSpace(input.Email))
	registrationName := strings.TrimSpace(input.Name)

	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		if _, err := store.FindUserByEmail(ctx, registrationEmail, false); err == nil {
			return NewAppError(http.StatusConflict, "email_taken", "email already exists", nil)
		} else if !repositories.IsNotFound(err) {
			return InternalError("failed to inspect existing user", err)
		}

		user, err := store.CreateUser(ctx, repositories.CreateUserParams{
			Name:         registrationName,
			Email:        registrationEmail,
			PasswordHash: passwordHash,
			Role:         "member",
			Status:       "active",
		})
		if err != nil {
			return InternalError("failed to create user", err)
		}

		issued, err := s.issueTokens(ctx, store, user)
		if err != nil {
			return err
		}

		if err := store.RecordUserLogin(ctx, user.ID, time.Now()); err != nil {
			return InternalError("failed to update last login", err)
		}

		result = issued
		return nil
	}); err != nil {
		return AuthResult{}, err
	}

	request.ActorUserID = &result.User.ID
	request.ActorEmail = result.User.Email
	request.ActorName = result.User.Name
	request.ActorRole = result.User.Role

	go createAuditLog(s.repo, request, "auth.register", "users", result.User.ID, map[string]any{
		"email": result.User.Email,
	})
	go s.sendRegistrationEmail(result.User)

	return result, nil
}

// REF.AUTH-02|Login
func (s *AuthService) Login(ctx context.Context, request RequestMetadata, input schemas.LoginRequest) (AuthResult, error) {
	user, err := s.repo.FindUserByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)), true)
	if err != nil {
		if repositories.IsNotFound(err) {
			return AuthResult{}, NewAppError(http.StatusUnauthorized, "invalid_credentials", "email or password is invalid", err)
		}

		return AuthResult{}, InternalError("failed to load user credentials", err)
	}

	if err := ensureUserAccessible(user); err != nil {
		return AuthResult{}, err
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

		if err := store.RecordUserLogin(ctx, user.ID, time.Now()); err != nil {
			return InternalError("failed to update last login", err)
		}

		result = issued
		return nil
	}); err != nil {
		return AuthResult{}, err
	}

	request.ActorUserID = &user.ID
	request.ActorName = user.Name
	request.ActorEmail = user.Email
	request.ActorRole = user.Role

	go createAuditLog(s.repo, request, "auth.login", "users", user.ID, map[string]any{"email": user.Email})

	return result, nil
}

// REF.AUTH-03|Logout
func (s *AuthService) Logout(ctx context.Context, request RequestMetadata, rawRefreshToken string) error {
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

	go createAuditLog(s.repo, request, "auth.logout", "sessions", valueOrEmpty(request.ActorUserID), nil)

	return nil
}

// REF.AUTH-04|Me
func (s *AuthService) Me(ctx context.Context, userID string) (repositories.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID, true)
	if err != nil {
		if repositories.IsNotFound(err) {
			return repositories.User{}, NewAppError(http.StatusUnauthorized, "unauthorized", "failed to resolve user", err)
		}

		return repositories.User{}, InternalError("failed to load authenticated user", err)
	}

	if err := ensureUserAccessible(user); err != nil {
		return repositories.User{}, err
	}

	return user, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, request RequestMetadata, input schemas.ForgotPasswordRequest) error {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	user, err := s.repo.FindUserByEmail(ctx, email, true)
	if err != nil {
		if repositories.IsNotFound(err) {
			go createAuditLog(s.repo, request, "auth.password-recovery.requested", "users", "", map[string]any{
				"email":   email,
				"outcome": "ignored",
			})
			return nil
		}

		return InternalError("failed to inspect user for password recovery", err)
	}

	if user.DeletedAt != nil || user.Status != "active" {
		go createAuditLog(s.repo, request, "auth.password-recovery.requested", "users", user.ID, map[string]any{
			"email":   email,
			"outcome": "ignored",
		})
		return nil
	}

	rawToken, err := generateOpaqueToken()
	if err != nil {
		return InternalError("failed to generate password reset token", err)
	}

	tokenHash := HashToken(rawToken)
	expiresAt := time.Now().Add(30 * time.Minute)

	if err := s.repo.SavePasswordResetToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return InternalError("failed to persist password reset token", err)
	}

	resetURL := s.frontendBaseURL + "/auth/reset-password?token=" + rawToken
	go createAuditLog(s.repo, request, "auth.password-recovery.requested", "users", user.ID, map[string]any{
		"email":     user.Email,
		"expiresAt": expiresAt.UTC(),
	})
	go s.sendPasswordResetEmail(user, resetURL)

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, request RequestMetadata, input schemas.ResetPasswordRequest) error {
	tokenHash := HashToken(strings.TrimSpace(input.Token))

	var resetUser repositories.User
	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		resetToken, err := store.FindPasswordResetToken(ctx, tokenHash)
		if err != nil {
			if repositories.IsNotFound(err) {
				return NewAppError(http.StatusBadRequest, "invalid_reset_token", "reset token is invalid", err)
			}
			return InternalError("failed to resolve password reset token", err)
		}

		if resetToken.UsedAt != nil || time.Now().After(resetToken.ExpiresAt) {
			return NewAppError(http.StatusBadRequest, "invalid_reset_token", "reset token is invalid or expired", nil)
		}

		user, err := store.FindUserByID(ctx, resetToken.UserID, true)
		if err != nil {
			return InternalError("failed to resolve user for password reset", err)
		}

		if err := ensureUserAccessible(user); err != nil {
			return err
		}

		passwordHash, err := HashPassword(input.NewPassword)
		if err != nil {
			return InternalError("failed to hash new password", err)
		}

		if _, err := store.UpdateUserPassword(ctx, user.ID, passwordHash); err != nil {
			return InternalError("failed to update password", err)
		}

		if err := store.UsePasswordResetToken(ctx, tokenHash); err != nil {
			return InternalError("failed to mark password reset token as used", err)
		}

		if err := store.RevokeRefreshTokensByUserID(ctx, user.ID); err != nil {
			return InternalError("failed to revoke refresh tokens after password reset", err)
		}

		resetUser = user
		return nil
	}); err != nil {
		return err
	}

	request.ActorUserID = &resetUser.ID
	request.ActorName = resetUser.Name
	request.ActorEmail = resetUser.Email
	request.ActorRole = resetUser.Role
	go createAuditLog(s.repo, request, "auth.password-reset.completed", "users", resetUser.ID, map[string]any{
		"email": resetUser.Email,
	})

	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, request RequestMetadata, userID string, input schemas.ChangePasswordRequest) error {
	user, err := s.repo.FindUserByID(ctx, userID, true)
	if err != nil {
		if repositories.IsNotFound(err) {
			return NewAppError(http.StatusUnauthorized, "unauthorized", "failed to resolve user", err)
		}
		return InternalError("failed to resolve current user", err)
	}

	if err := ensureUserAccessible(user); err != nil {
		return err
	}

	valid, err := VerifyPassword(user.PasswordHash, input.CurrentPassword)
	if err != nil {
		return InternalError("failed to verify current password", err)
	}

	if !valid {
		return NewAppError(http.StatusBadRequest, "invalid_password", "current password is invalid", nil)
	}

	passwordHash, err := HashPassword(input.NewPassword)
	if err != nil {
		return InternalError("failed to hash new password", err)
	}

	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		if _, err := store.UpdateUserPassword(ctx, user.ID, passwordHash); err != nil {
			return InternalError("failed to update password", err)
		}

		if err := store.RevokeRefreshTokensByUserID(ctx, user.ID); err != nil {
			return InternalError("failed to revoke refresh tokens after password change", err)
		}

		return nil
	}); err != nil {
		return err
	}

	request.ActorUserID = &user.ID
	request.ActorName = user.Name
	request.ActorEmail = user.Email
	request.ActorRole = user.Role
	go createAuditLog(s.repo, request, "auth.password-changed", "users", user.ID, map[string]any{
		"email": user.Email,
	})

	return nil
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

func (s *AuthService) sendPasswordResetEmail(user repositories.User, resetURL string) {
	if s.email == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.email.SendPasswordResetNotice(ctx, user, resetURL); err != nil {
		log.Printf("failed to send password reset email to %s: %v", user.Email, err)
	}
}

func ensureUserAccessible(user repositories.User) error {
	if user.DeletedAt != nil {
		return NewAppError(http.StatusForbidden, "account_deleted", "account is deleted", nil)
	}
	if user.Status != "active" {
		return NewAppError(http.StatusForbidden, "account_inactive", "account is inactive", nil)
	}
	return nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func generateOpaqueToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
