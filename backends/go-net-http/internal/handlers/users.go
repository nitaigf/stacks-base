package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
	"stacks-base/backends/go-net-http/internal/services"
	"stacks-base/backends/go-net-http/internal/utils"
)

type UsersHandler struct {
	service *services.UserService
}

func NewUsersHandler(service *services.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	query := schemas.ParseUserListQuery(readQueryMap(r.URL.Query()))
	users, meta, err := h.service.ListUsers(r.Context(), requestMetadataFromRequest(r), query)
	if err != nil {
		writeServiceError(w, err, "failed to list users")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": sanitizeUsers(users),
		"meta": meta,
	})
}

func (h *UsersHandler) Show(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(r.Context(), requestMetadataFromRequest(r), r.PathValue("userID"))
	if err != nil {
		writeServiceError(w, err, "failed to load user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"data": sanitizeUser(user)})
}

func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload schemas.UserCreateRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "request body is invalid", nil)
		return
	}

	if validationErrors := payload.Validate(); len(validationErrors) > 0 {
		utils.WriteError(w, http.StatusBadRequest, "validation_error", "request body failed validation", validationErrors)
		return
	}

	user, err := h.service.CreateUser(r.Context(), requestMetadataFromRequest(r), payload)
	if err != nil {
		writeServiceError(w, err, "failed to create user")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{"data": sanitizeUser(user)})
}

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	var payload schemas.UserUpdateRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "request body is invalid", nil)
		return
	}

	if validationErrors := payload.Validate(); len(validationErrors) > 0 {
		utils.WriteError(w, http.StatusBadRequest, "validation_error", "request body failed validation", validationErrors)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), requestMetadataFromRequest(r), r.PathValue("userID"), payload)
	if err != nil {
		writeServiceError(w, err, "failed to update user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"data": sanitizeUser(user)})
}

func (h *UsersHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	h.updateUserAction(w, r, h.service.DeactivateUser, "failed to deactivate user")
}

func (h *UsersHandler) Reactivate(w http.ResponseWriter, r *http.Request) {
	h.updateUserAction(w, r, h.service.ReactivateUser, "failed to reactivate user")
}

func (h *UsersHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	h.updateUserAction(w, r, h.service.SoftDeleteUser, "failed to soft-delete user")
}

func (h *UsersHandler) Restore(w http.ResponseWriter, r *http.Request) {
	h.updateUserAction(w, r, h.service.RestoreUser, "failed to restore user")
}

func (h *UsersHandler) HardDelete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.HardDeleteUser(r.Context(), requestMetadataFromRequest(r), r.PathValue("userID")); err != nil {
		writeServiceError(w, err, "failed to hard-delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UsersHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	query := schemas.ParseUserListQuery(readQueryMap(r.URL.Query()))
	users, err := h.service.ExportUsers(r.Context(), requestMetadataFromRequest(r), query, "users.export.csv")
	if err != nil {
		writeServiceError(w, err, "failed to export users as csv")
		return
	}

	payload, err := utils.BuildUsersCSV(users)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to build csv export", nil)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="users.csv"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (h *UsersHandler) ExportXLSX(w http.ResponseWriter, r *http.Request) {
	query := schemas.ParseUserListQuery(readQueryMap(r.URL.Query()))
	users, err := h.service.ExportUsers(r.Context(), requestMetadataFromRequest(r), query, "users.export.xlsx")
	if err != nil {
		writeServiceError(w, err, "failed to export users as xlsx")
		return
	}

	payload, err := utils.BuildUsersXLSX(users)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to build xlsx export", nil)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="users.xlsx"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (h *UsersHandler) Print(w http.ResponseWriter, r *http.Request) {
	query := schemas.ParseUserListQuery(readQueryMap(r.URL.Query()))
	users, err := h.service.ExportUsers(r.Context(), requestMetadataFromRequest(r), query, "users.print")
	if err != nil {
		writeServiceError(w, err, "failed to render printable users report")
		return
	}

	document, err := utils.BuildUsersPDF(users)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to render printable report", nil)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `inline; filename="users.pdf"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(document)
}

func (h *UsersHandler) updateUserAction(
	w http.ResponseWriter,
	r *http.Request,
	fn func(ctx context.Context, request services.RequestMetadata, id string) (repositories.User, error),
	fallbackMessage string,
) {
	user, err := fn(r.Context(), requestMetadataFromRequest(r), r.PathValue("userID"))
	if err != nil {
		writeServiceError(w, err, fallbackMessage)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"data": sanitizeUser(user)})
}

func sanitizeUsers(users []repositories.User) []map[string]any {
	payload := make([]map[string]any, 0, len(users))
	for _, user := range users {
		payload = append(payload, sanitizeUser(user))
	}
	return payload
}

func readQueryMap(values url.Values) map[string]string {
	result := make(map[string]string, len(values))
	for key, value := range values {
		if len(value) == 0 {
			continue
		}
		result[key] = value[0]
	}
	return result
}

func writeServiceError(w http.ResponseWriter, err error, fallbackMessage string) {
	var appError *services.AppError
	if errors.As(err, &appError) {
		utils.WriteError(w, appError.StatusCode, appError.Code, appError.Message, nil)
		return
	}

	utils.WriteError(w, http.StatusInternalServerError, "internal_error", fallbackMessage, nil)
}
