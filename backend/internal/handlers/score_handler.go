package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ScoreHandler struct {
	service *services.ScoreService
}

func NewScoreHandler(s *services.ScoreService) *ScoreHandler {
	return &ScoreHandler{service: s}
}

// Depending on the request parameters, this handler returns scores for a user_id, a forecast_id, or both
func (h *ScoreHandler) GetScores(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	var userID int64
	userIDstr := queryParams.Get("user_id")
	if userIDstr != "" {
		var err error
		userID, err = strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
	}

	var forecastID int64
	forecastIDstr := queryParams.Get("forecast_id")
	if forecastIDstr != "" {
		var err error
		forecastID, err = strconv.ParseInt(forecastIDstr, 10, 64)
		if err != nil {
			http.Error(w, "invalid forecast ID", http.StatusBadRequest)
			return
		}
	}

	scores, err := h.service.GetScores(r.Context(), userID, forecastID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)

}

// Handlers to modify/create scores
func (h *ScoreHandler) CreateScore(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var score models.Scores
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure score is created for authenticated user
	score.UserID = claims.UserID

	if err := h.service.CreateScore(r.Context(), &score); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, score.ID)
}

func (h *ScoreHandler) DeleteScore(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var score models.Scores
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify user owns this score
	if score.UserID != claims.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteScore(r.Context(), score.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, "Score deleted successfully")
}

func (h *ScoreHandler) GetAverageScores(w http.ResponseWriter, r *http.Request) {
	scores, err := h.service.GetAverageScores(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScores(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	var userIDPtr *int64
	userIDstr := queryParams.Get("user_id")
	if userIDstr != "" {
		userID, err := strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		userIDPtr = &userID
	}

	var forecastIDPtr *int64
	forecastIDstr := queryParams.Get("forecast_id")
	if forecastIDstr != "" {
		forecastID, err := strconv.ParseInt(forecastIDstr, 10, 64)
		if err != nil {
			http.Error(w, "invalid forecast ID", http.StatusBadRequest)
			return
		}
		forecastIDPtr = &forecastID
	}

	var categoryPtr *string
	categorystr := queryParams.Get("category")
	if categorystr != "" {
		category := strings.ToLower(categorystr)
		categoryPtr = &category
	}

	scores, err := h.service.GetAggregateScores(r.Context(), userIDPtr, forecastIDPtr, categoryPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresGroupedByUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	var categoryPtr *string
	categorystr := queryParams.Get("category")
	if categorystr != "" {
		category := strings.ToLower(categorystr)
		categoryPtr = &category
	}

	scores, err := h.service.GetAggregateScoresGroupedByUsers(r.Context(), categoryPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}
