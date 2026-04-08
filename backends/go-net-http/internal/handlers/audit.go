package handlers

import (
	"net/http"

	"stacks-base/backends/go-net-http/internal/schemas"
	"stacks-base/backends/go-net-http/internal/services"
	"stacks-base/backends/go-net-http/internal/utils"
)

type AuditHandler struct {
	service *services.AuditService
}

func NewAuditHandler(service *services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	query := schemas.ParseAuditLogListQuery(readQueryMap(r.URL.Query()))
	logs, meta, err := h.service.ListAuditLogs(r.Context(), requestMetadataFromRequest(r), query)
	if err != nil {
		writeServiceError(w, err, "failed to list audit logs")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": logs,
		"meta": meta,
	})
}
