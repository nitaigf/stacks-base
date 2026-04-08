package handlers

import (
	"net"
	"net/http"
	"strings"

	"stacks-base/backends/go-net-http/internal/middleware"
	"stacks-base/backends/go-net-http/internal/services"
)

func requestMetadataFromRequest(r *http.Request) services.RequestMetadata {
	metadata := services.RequestMetadata{
		Route:     r.URL.Path,
		Method:    r.Method,
		IPAddress: clientIP(r),
		UserAgent: strings.TrimSpace(r.UserAgent()),
	}

	if claims, ok := middleware.AuthClaimsFromContext(r.Context()); ok {
		metadata.ActorUserID = &claims.UserID
		metadata.ActorEmail = claims.Email
		metadata.ActorRole = claims.Role
	}

	return metadata
}

func clientIP(r *http.Request) string {
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	forwarded := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if forwarded != "" {
		return forwarded
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}

	return host
}
