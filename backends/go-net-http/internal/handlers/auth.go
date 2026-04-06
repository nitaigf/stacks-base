package handlers

import (
	"errors"
	"net/http"
	"time"

	"stacks-base/backends/go-net-http/internal/middleware"
	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
	"stacks-base/backends/go-net-http/internal/services"
	"stacks-base/backends/go-net-http/internal/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// REF.AUTH-01|Register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload schemas.RegisterRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "request body is invalid", nil)
		return
	}

	if validationErrors := payload.Validate(); len(validationErrors) > 0 {
		utils.WriteError(w, http.StatusBadRequest, "validation_error", "request body failed validation", validationErrors)
		return
	}

	result, err := h.service.Register(r.Context(), payload)
	if err != nil {
		h.writeServiceError(w, err, "failed to register user")
		return
	}

	h.setRefreshCookie(w, result.RefreshToken, result.RefreshExpiresAt)
	utils.WriteJSON(w, http.StatusCreated, authEnvelope(result))
}

// REF.AUTH-02|Login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload schemas.LoginRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "request body is invalid", nil)
		return
	}

	if validationErrors := payload.Validate(); len(validationErrors) > 0 {
		utils.WriteError(w, http.StatusBadRequest, "validation_error", "request body failed validation", validationErrors)
		return
	}

	result, err := h.service.Login(r.Context(), payload)
	if err != nil {
		h.writeServiceError(w, err, "failed to authenticate user")
		return
	}

	h.setRefreshCookie(w, result.RefreshToken, result.RefreshExpiresAt)
	utils.WriteJSON(w, http.StatusOK, authEnvelope(result))
}

// REF.AUTH-03|Logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.AuthClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "authentication is required", nil)
		return
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "refresh cookie is required", nil)
		return
	}

	if err := h.service.Logout(r.Context(), refreshCookie.Value); err != nil {
		h.writeServiceError(w, err, "failed to revoke session")
		return
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// REF.AUTH-04|Me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.AuthClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "authentication is required", nil)
		return
	}

	user, err := h.service.Me(r.Context(), claims.UserID)
	if err != nil {
		h.writeServiceError(w, err, "failed to resolve user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"data": sanitizeUser(user)})
}

func authEnvelope(result services.AuthResult) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"accessToken": result.AccessToken,
			"user":        sanitizeUser(result.User),
		},
	}
}

func sanitizeUser(user repositories.User) map[string]any {
	return map[string]any{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"status":    user.Status,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func (h *AuthHandler) writeServiceError(w http.ResponseWriter, err error, fallbackMessage string) {
	var appError *services.AppError
	if errors.As(err, &appError) {
		utils.WriteError(w, appError.StatusCode, appError.Code, appError.Message, nil)
		return
	}

	utils.WriteError(w, http.StatusInternalServerError, "internal_error", fallbackMessage, nil)
}