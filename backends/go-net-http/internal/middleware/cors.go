package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

func AllowCORS(origin string, next http.Handler) http.Handler {
	allowedOrigins := buildAllowedOrigins(origin)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestOrigin := strings.TrimSpace(r.Header.Get("Origin"))
		if requestOrigin != "" && allowedOrigins[requestOrigin] {
			w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func buildAllowedOrigins(configuredOrigin string) map[string]bool {
	allowed := map[string]bool{}
	configuredOrigin = strings.TrimSpace(configuredOrigin)
	if configuredOrigin == "" {
		return allowed
	}

	allowed[configuredOrigin] = true

	parsed, err := url.Parse(configuredOrigin)
	if err != nil || parsed.Host == "" {
		return allowed
	}

	hostname := parsed.Hostname()
	port := parsed.Port()
	scheme := parsed.Scheme

	if hostname == "127.0.0.1" {
		host := "localhost"
		if port != "" {
			host = host + ":" + port
		}
		allowed[scheme+"://"+host] = true
	}

	if hostname == "localhost" {
		host := "127.0.0.1"
		if port != "" {
			host = host + ":" + port
		}
		allowed[scheme+"://"+host] = true
	}

	return allowed
}