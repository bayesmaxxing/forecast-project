package handlers

import (
	"backend/internal/logger"
	"backend/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"
)

type NewsHandler struct {
	service *services.NewsService
}

func NewNewsHandler(s *services.NewsService) *NewsHandler {
	return &NewsHandler{service: s}
}

func (h *NewsHandler) GetNews(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	var request struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error("invalid request body", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.Query == "" {
		log.Error("query is required")
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	log.Info("getting news", slog.String("query", request.Query))
	news, err := h.service.GetNews(r.Context(), request.Query)
	if err != nil {
		log.Error("error getting news", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newsContent := news.GetLastContent()

	log.Info("news retrieved", slog.String("news_content", newsContent))
	respondJSON(w, http.StatusOK, newsContent)
}
