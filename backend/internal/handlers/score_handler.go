package handlers

import (
	"backend/internal/auth"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ScoreHandler struct {
	service *services.ScoreService
}

func NewScoreHandler(s *services.ScoreService) *ScoreHandler {
	return &ScoreHandler{service: s}
}

// Depending on the request parameters, this handler returns scores for a user_id, a forecast_id, or both
func (h *ScoreHandler) GetScores(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	queryParams := r.URL.Query()

	var userID int64
	userIDstr := queryParams.Get("user_id")
	if userIDstr != "" {
		var err error
		userID, err = strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			log.Error("invalid user ID", slog.String("error", err.Error()))
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
			log.Error("invalid forecast ID", slog.String("error", err.Error()))
			http.Error(w, "invalid forecast ID", http.StatusBadRequest)
			return
		}
	}

	log.Info("getting scores", slog.Int64("user_id", userID), slog.Int64("forecast_id", forecastID))
	scores, err := h.service.GetScores(r.Context(), userID, forecastID)
	if err != nil {
		log.Error("failed to get scores", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)

}

// Handlers to modify/create scores
func (h *ScoreHandler) CreateScore(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		log.Error("unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var score models.Scores
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		log.Error("invalid request body", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure score is created for authenticated user
	score.UserID = claims.UserID

	log.Info("creating score", slog.Any("score", score))
	if err := h.service.CreateScore(r.Context(), &score); err != nil {
		log.Error("failed to create score", slog.String("error", err.Error()))
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
	log := logger.FromContext(r.Context())
	log.Info("getting average scores")
	scores, err := h.service.GetAverageScores(r.Context())
	if err != nil {
		log.Error("failed to get average scores", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScores(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	queryParams := r.URL.Query()

	var userIDPtr *int64
	userIDstr := queryParams.Get("user_id")
	if userIDstr != "" {
		userID, err := strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			log.Error("invalid user ID", slog.String("error", err.Error()))
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
			log.Error("invalid forecast ID", slog.String("error", err.Error()))
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

	var startDatePtr *time.Time
	startDateStr := queryParams.Get("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			log.Error("invalid start_date", slog.String("error", err.Error()))
			http.Error(w, "invalid start_date, expected RFC3339 format", http.StatusBadRequest)
			return
		}
		startDatePtr = &startDate
	}

	var endDatePtr *time.Time
	endDateStr := queryParams.Get("end_date")
	if endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			log.Error("invalid end_date", slog.String("error", err.Error()))
			http.Error(w, "invalid end_date, expected RFC3339 format", http.StatusBadRequest)
			return
		}
		endDatePtr = &endDate
	}

	log.Info("getting aggregate scores", slog.Any("user_id", userIDPtr), slog.Any("forecast_id", forecastIDPtr), slog.Any("category", categoryPtr), slog.Any("start_date", startDatePtr), slog.Any("end_date", endDatePtr))
	scores, err := h.service.GetAggregateScores(r.Context(), userIDPtr, forecastIDPtr, categoryPtr, startDatePtr, endDatePtr)
	if err != nil {
		log.Error("failed to get aggregate scores", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}

func (h *ScoreHandler) GetAggregateScoresGroupedByUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	queryParams := r.URL.Query()

	var categoryPtr *string
	categorystr := queryParams.Get("category")
	if categorystr != "" {
		category := strings.ToLower(categorystr)
		categoryPtr = &category
	}

	var startDatePtr *time.Time
	startDateStr := queryParams.Get("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			log.Error("invalid start_date", slog.String("error", err.Error()))
			http.Error(w, "invalid start_date, expected RFC3339 format", http.StatusBadRequest)
			return
		}
		startDatePtr = &startDate
	}

	var endDatePtr *time.Time
	endDateStr := queryParams.Get("end_date")
	if endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			log.Error("invalid end_date", slog.String("error", err.Error()))
			http.Error(w, "invalid end_date, expected RFC3339 format", http.StatusBadRequest)
			return
		}
		endDatePtr = &endDate
	}

	log.Info("getting aggregate scores grouped by users", slog.Any("category", categoryPtr), slog.Any("start_date", startDatePtr), slog.Any("end_date", endDatePtr))
	scores, err := h.service.GetAggregateScoresGroupedByUsers(r.Context(), categoryPtr, startDatePtr, endDatePtr)
	if err != nil {
		log.Error("failed to get aggregate scores grouped by users", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, scores)
}
