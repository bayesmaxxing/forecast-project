package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type ScoreHandler struct {
	service *services.ScoreService
}

func NewScoreHandler(s *services.ScoreService) *ScoreHandler {
	return &ScoreHandler{service: s}
}

// Depending on the request parameters, this handler returns scores for a user_id, a forecast_id, or both
func (h *ScoreHandler) GetScores(w http.ResponseWriter, r *http.Request) {
	var request struct {
		User_id     *int64 `json:"user_id"`
		Forecast_id *int64 `json:"forecast_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := *request.User_id
	forecastID := *request.Forecast_id

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

func (h *ScoreHandler) GetAllScores(w http.ResponseWriter, r *http.Request) {

	scores, err := h.service.GetAllScores(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAverageScores(w http.ResponseWriter, r *http.Request) {
	scores, err := h.service.GetAverageScores(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, scores)
}

// Gets the average score for a specific forecast
func (h *ScoreHandler) GetAverageScoreByForecastID(w http.ResponseWriter, r *http.Request) {
	forecastID := r.PathValue("id")
	if forecastID == "" {
		http.Error(w, "forecast ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(forecastID, 10, 64)
	if err != nil {
		http.Error(w, "invalid forecast ID", http.StatusBadRequest)
		return
	}

	avgScore, err := h.service.GetAverageScoreByForecastID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, avgScore)
}

func (h *ScoreHandler) GetAggregateScores(w http.ResponseWriter, r *http.Request) {
	scores, err := h.service.GetOverallScores(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	scores, err := h.service.GetAggregateScoresByUserID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresByUserIDAndCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}
	category := r.PathValue("category")
	if category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	scores, err := h.service.GetAggregateScoresByUserIDAndCategory(r.Context(), id, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresByUsers(w http.ResponseWriter, r *http.Request) {
	scores, err := h.service.GetAggregateScoresByUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}

	scores, err := h.service.GetAggregateScoresByCategory(r.Context(), category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresByUsersAndCategory(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}
	scores, err := h.service.GetAggregateScoresByUsersAndCategory(r.Context(), category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}
