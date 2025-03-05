package handlers

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type ScoreHandler struct {
	service *services.ScoreService
	cache   *cache.Cache
}

func NewScoreHandler(s *services.ScoreService, c *cache.Cache) *ScoreHandler {
	return &ScoreHandler{service: s, cache: c}
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

	switch {
	case request.User_id != nil && request.Forecast_id != nil:
		scores, err := h.service.GetScoreByForecastAndUser(r.Context(), *request.Forecast_id, *request.User_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, scores)
	case request.User_id != nil:
		scores, err := h.service.GetScoresByUserID(r.Context(), *request.User_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, scores)
	case request.Forecast_id != nil:
		scores, err := h.service.GetScoreByForecastID(r.Context(), *request.Forecast_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, scores)
	case request.User_id == nil && request.Forecast_id == nil:
		scores, err := h.service.GetAllScores(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, scores)
	default:
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
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

	// Generate cache key
	cacheKey := "avg_score_forecast_" + forecastID

	// Try to get from cache first
	if cachedData, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedData)
		return
	}

	avgScore, err := h.service.GetAverageScoreByForecastID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in cache
	h.cache.Set(cacheKey, avgScore)
	respondJSON(w, http.StatusOK, avgScore)
}

// Handler for aggregate scores
type aggregateScoresRequest struct {
	Category *string `json:"category"`
	User_id  *int64  `json:"user_id"`
	ByUser   *bool   `json:"by_user"`
}

func (h *ScoreHandler) GetAggregateScores(w http.ResponseWriter, r *http.Request) {
	var request aggregateScoresRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set default value for ByUser if it's nil
	byUser := false
	if request.ByUser != nil {
		byUser = *request.ByUser
	}

	// Generate cache key based on request parameters
	cacheKey := "aggregate_scores"
	if request.User_id != nil {
		cacheKey += "_user_" + strconv.FormatInt(*request.User_id, 10)
	}
	if request.Category != nil {
		cacheKey += "_cat_" + *request.Category
	}
	if byUser && request.User_id == nil {
		cacheKey += "_by_user"
	}

	// Try to get from cache first
	if cachedData, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedData)
		return
	}

	var scores any
	var err error

	switch {
	// For specific user cases
	case request.User_id != nil && request.Category != nil:
		scores, err = h.service.GetUserCategoryScores(r.Context(), *request.User_id, *request.Category)
	case request.User_id != nil && request.Category == nil:
		scores, err = h.service.GetUserOverallScores(r.Context(), *request.User_id)

	// For all users cases
	case request.User_id == nil && request.Category != nil && byUser:
		scores, err = h.service.GetCategoryScoresByUsers(r.Context(), *request.Category)
	case request.User_id == nil && request.Category == nil && byUser:
		scores, err = h.service.GetOverallScoresByUsers(r.Context())
	case request.User_id == nil && request.Category != nil && !byUser:
		scores, err = h.service.GetCategoryScores(r.Context(), *request.Category)
	case request.User_id == nil && request.Category == nil && !byUser:
		scores, err = h.service.GetOverallScores(r.Context())
	default:
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in cache
	h.cache.Set(cacheKey, scores)
	respondJSON(w, http.StatusOK, scores)
}
