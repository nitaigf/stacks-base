package routes

import (
	"net/http"

	"stacks-base/backends/go-net-http/internal/config"
	"stacks-base/backends/go-net-http/internal/handlers"
	"stacks-base/backends/go-net-http/internal/middleware"
	"stacks-base/backends/go-net-http/internal/services"
)

func New(
	cfg config.Config,
	authService *services.AuthService,
	userService *services.UserService,
	auditService *services.AuditService,
) http.Handler {
	mux := http.NewServeMux()
	healthHandler := handlers.HealthHandler{}
	authHandler := handlers.NewAuthHandler(authService)
	usersHandler := handlers.NewUsersHandler(userService)
	auditHandler := handlers.NewAuditHandler(auditService)

	mux.Handle("GET /health", healthHandler)
	// REF.AUTH-01|Register
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	// REF.AUTH-02|Login
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)
	// REF.AUTH-03|Logout
	mux.Handle("POST /api/v1/auth/logout", middleware.RequireAuth(authService, http.HandlerFunc(authHandler.Logout)))
	// REF.AUTH-04|Me
	mux.Handle("GET /api/v1/users/me", middleware.RequireAuth(authService, http.HandlerFunc(authHandler.Me)))
	mux.Handle("POST /api/v1/auth/change-password", middleware.RequireAuth(authService, http.HandlerFunc(authHandler.ChangePassword)))
	mux.Handle("GET /api/v1/users", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.List)))
	mux.Handle("POST /api/v1/users", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Create)))
	mux.Handle("GET /api/v1/users/export.csv", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.ExportCSV)))
	mux.Handle("GET /api/v1/users/export.xlsx", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.ExportXLSX)))
	mux.Handle("GET /api/v1/users/print", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Print)))
	mux.Handle("GET /api/v1/users/{userID}", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Show)))
	mux.Handle("PATCH /api/v1/users/{userID}", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Update)))
	mux.Handle("POST /api/v1/users/{userID}/deactivate", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Deactivate)))
	mux.Handle("POST /api/v1/users/{userID}/reactivate", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Reactivate)))
	mux.Handle("POST /api/v1/users/{userID}/soft-delete", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.SoftDelete)))
	mux.Handle("POST /api/v1/users/{userID}/restore", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.Restore)))
	mux.Handle("DELETE /api/v1/users/{userID}", middleware.RequireAuth(authService, http.HandlerFunc(usersHandler.HardDelete)))
	mux.Handle("GET /api/v1/audit-logs", middleware.RequireAuth(authService, http.HandlerFunc(auditHandler.List)))

	return middleware.AllowCORS(cfg.AllowedOrigin, mux)
}
