package handlers

import (
	"backend/internal/services"
	"encoding/json"
	"net/http"
)

type NewsHandler struct {
	service *services.NewsService
}

func NewNewsHandler(s *services.NewsService) *NewsHandler {
	return &NewsHandler{service: s}
}

func (h *NewsHandler) GetNews(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.Query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	news, err := h.service.GetNews(r.Context(), request.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newsContent := news.GetLastContent()

	respondJSON(w, http.StatusOK, newsContent)
}
