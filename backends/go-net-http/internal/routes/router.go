package routes

import (
	"net/http"

	"stacks-base/backends/go-net-http/internal/config"
	"stacks-base/backends/go-net-http/internal/handlers"
	"stacks-base/backends/go-net-http/internal/middleware"
	"stacks-base/backends/go-net-http/internal/services"
)

func New(cfg config.Config, authService *services.AuthService) http.Handler {
	mux := http.NewServeMux()
	healthHandler := handlers.HealthHandler{}
	authHandler := handlers.NewAuthHandler(authService)

	mux.Handle("GET /health", healthHandler)
	// REF.AUTH-01|Register
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	// REF.AUTH-02|Login
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	// REF.AUTH-03|Logout
	mux.Handle("POST /api/v1/auth/logout", middleware.RequireAuth(authService, http.HandlerFunc(authHandler.Logout)))
	// REF.AUTH-04|Me
	mux.Handle("GET /api/v1/users/me", middleware.RequireAuth(authService, http.HandlerFunc(authHandler.Me)))

	return middleware.AllowCORS(cfg.AllowedOrigin, mux)
}