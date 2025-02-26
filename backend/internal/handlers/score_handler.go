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

// Handlers for aggregate scores
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

	if request.User_id != nil {
		http.Error(w, "user_id is not allowed for this endpoint", http.StatusBadRequest)
		return
	}

	// Generate cache key based on request parameters
	cacheKey := "aggregate_scores"
	if request.Category != nil {
		cacheKey += "_cat_" + *request.Category
	}
	if request.ByUser != nil && *request.ByUser {
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
	case request.Category != nil && *request.ByUser && request.User_id == nil:
		scores, err = h.service.GetCategoryScoresByUsers(r.Context(), *request.Category)
	case request.Category == nil && *request.ByUser && request.User_id == nil:
		scores, err = h.service.GetOverallScoresByUsers(r.Context())
	case request.Category != nil && !*request.ByUser && request.User_id == nil:
		scores, err = h.service.GetCategoryScores(r.Context(), *request.Category)
	case request.Category == nil && !*request.ByUser && request.User_id == nil:
		scores, err = h.service.GetOverallScores(r.Context())
	default:
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in cache (with default expiration time)
	h.cache.Set(cacheKey, scores)
	respondJSON(w, http.StatusOK, scores)
}

// Handler for user-specific aggregate scores
func (h *ScoreHandler) GetUserAggregateScores(w http.ResponseWriter, r *http.Request) {
	var request aggregateScoresRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.User_id == nil {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if request.ByUser != nil {
		http.Error(w, "by_user is not allowed for this endpoint", http.StatusBadRequest)
		return
	}

	// Generate cache key based on request parameters
	cacheKey := "user_aggregate_scores_" + strconv.FormatInt(*request.User_id, 10)
	if request.Category != nil {
		cacheKey += "_cat_" + *request.Category
	}

	// Try to get from cache first
	if cachedData, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedData)
		return
	}

	var scores any
	var err error

	switch {
	case request.Category != nil:
		scores, err = h.service.GetUserCategoryScores(r.Context(), *request.User_id, *request.Category)
	case request.Category == nil:
		scores, err = h.service.GetUserOverallScores(r.Context(), *request.User_id)
	default:
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, scores)
	respondJSON(w, http.StatusOK, scores)
}
