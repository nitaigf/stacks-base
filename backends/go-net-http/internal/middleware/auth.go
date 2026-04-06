package middleware

import (
	"context"
	"net/http"
	"strings"

	"stacks-base/backends/go-net-http/internal/services"
	"stacks-base/backends/go-net-http/internal/utils"
)

type contextKey string

const authClaimsContextKey contextKey = "authClaims"

type AccessTokenParser interface {
	ParseAccessToken(raw string) (services.AuthClaims, error)
}

func RequireAuth(parser AccessTokenParser, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(authorization, "Bearer ") {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing bearer token", nil)
			return
		}

		rawToken := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		claims, err := parser.ParseAccessToken(rawToken)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid access token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), authClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthClaimsFromContext(ctx context.Context) (services.AuthClaims, bool) {
	claims, ok := ctx.Value(authClaimsContextKey).(services.AuthClaims)
	return claims, ok
}