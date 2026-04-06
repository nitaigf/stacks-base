package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	paseto "aidanwoods.dev/go-paseto"

	"stacks-base/backends/go-net-http/internal/config"
	"stacks-base/backends/go-net-http/internal/repositories"
)

type AuthClaims struct {
	UserID string
	Email  string
	Role   string
}

type TokenService struct {
	accessKey  paseto.V4SymmetricKey
	refreshKey paseto.V4SymmetricKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(cfg config.Config) (*TokenService, error) {
	accessKey, err := paseto.V4SymmetricKeyFromBytes(hashKey(cfg.AccessTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("build access key: %w", err)
	}

	refreshKey, err := paseto.V4SymmetricKeyFromBytes(hashKey(cfg.RefreshTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("build refresh key: %w", err)
	}

	return &TokenService{
		accessKey:  accessKey,
		refreshKey: refreshKey,
		accessTTL:  cfg.AccessTokenTTL,
		refreshTTL: cfg.RefreshTokenTTL,
	}, nil
}

func (s *TokenService) IssueTokens(user repositories.User) (string, string, time.Time, error) {
	now := time.Now()
	validFrom := now.Add(-5 * time.Second)
	accessExpiration := now.Add(s.accessTTL)
	refreshExpiration := now.Add(s.refreshTTL)

	access := paseto.NewToken()
	access.SetIssuedAt(validFrom)
	access.SetNotBefore(validFrom)
	access.SetExpiration(accessExpiration)
	access.SetSubject(user.ID)
	access.SetIssuer("stacks-base")
	access.SetString("userId", user.ID)
	access.SetString("email", user.Email)
	access.SetString("role", user.Role)

	refresh := paseto.NewToken()
	refresh.SetIssuedAt(validFrom)
	refresh.SetNotBefore(validFrom)
	refresh.SetExpiration(refreshExpiration)
	refresh.SetSubject(user.ID)
	refresh.SetIssuer("stacks-base")
	refresh.SetString("userId", user.ID)
	refresh.SetString("email", user.Email)
	refresh.SetString("role", user.Role)
	refresh.SetString("type", "refresh")

	return access.V4Encrypt(s.accessKey, nil), refresh.V4Encrypt(s.refreshKey, nil), refreshExpiration, nil
}

func (s *TokenService) ParseAccessToken(raw string) (AuthClaims, error) {
	token, err := currentParser().ParseV4Local(s.accessKey, raw, nil)
	if err != nil {
		return AuthClaims{}, fmt.Errorf("parse access token: %w", err)
	}

	return claimsFromToken(token)
}

func (s *TokenService) ParseRefreshToken(raw string) (AuthClaims, error) {
	token, err := currentParser().ParseV4Local(s.refreshKey, raw, nil)
	if err != nil {
		return AuthClaims{}, fmt.Errorf("parse refresh token: %w", err)
	}

	claimType, err := token.GetString("type")
	if err != nil || claimType != "refresh" {
		return AuthClaims{}, fmt.Errorf("invalid refresh token type")
	}

	return claimsFromToken(token)
}

func HashToken(raw string) string {
	digest := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(digest[:])
}

func claimsFromToken(token *paseto.Token) (AuthClaims, error) {
	userID, err := token.GetString("userId")
	if err != nil {
		return AuthClaims{}, fmt.Errorf("read userId claim: %w", err)
	}

	email, err := token.GetString("email")
	if err != nil {
		return AuthClaims{}, fmt.Errorf("read email claim: %w", err)
	}

	role, err := token.GetString("role")
	if err != nil {
		return AuthClaims{}, fmt.Errorf("read role claim: %w", err)
	}

	return AuthClaims{UserID: userID, Email: email, Role: role}, nil
}

func hashKey(secret string) []byte {
	digest := sha256.Sum256([]byte(secret))
	return digest[:]
}

func currentParser() paseto.Parser {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	return parser
}