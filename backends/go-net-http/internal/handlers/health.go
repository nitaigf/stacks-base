package handlers

import (
	"net/http"

	"stacks-base/backends/go-net-http/internal/utils"
)

type HealthHandler struct{}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}