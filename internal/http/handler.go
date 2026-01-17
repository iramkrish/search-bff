package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/iramkrish/search-bff/internal/search"
)

type Handler struct {
	logger  *log.Logger
	service *search.Service
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger:  logger,
		service: search.NewService(logger),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "missing query param: q")
		return
	}

	resp, err := h.service.Search(r.Context(), query)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, search.ErrUpstreamUnavailable) {
			status = http.StatusBadGateway
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
