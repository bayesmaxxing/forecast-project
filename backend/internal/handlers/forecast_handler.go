package handlers

import (
	"backend/internal/auth"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
)

type ForecastHandler struct {
	service *services.ForecastService
}

func NewForecastHandler(s *services.ForecastService) *ForecastHandler {
	return &ForecastHandler{service: s}
}

// handlers for lists of forecasts
func (h *ForecastHandler) ListForecasts(w http.ResponseWriter, r *http.Request) {

	log := logger.FromContext(r.Context())

	queryParams := r.URL.Query()

	// ensure forecast_id is not part of the query params
	forecastIDstr := queryParams.Get("forecast_id")
	if forecastIDstr != "" {
		forecastID, err := strconv.ParseInt(forecastIDstr, 10, 64)
		if err != nil {
			log.Warn("invalid forecast_id format", slog.String("error", err.Error()))
			http.Error(w, "invalid forecast_id format", http.StatusBadRequest)
			return
		}

		if forecastID == 0 {
			log.Warn("forecast_id is not supported for this operation")
			http.Error(w, "forecast_id is not supported for this operation", http.StatusBadRequest)
			return
		}
	}

	filters := models.ForecastFilters{}

	status := queryParams.Get("status")
	if status != "" {
		// validate status
		validStatuses := []string{"open", "resolved", "closed"}
		if !slices.Contains(validStatuses, status) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		filters.Status = &status
	}

	category := queryParams.Get("category")
	if category != "" {
		filters.Category = &category
	}

	log.Info("getting forecasts", slog.Any("filters", filters))
	forecasts, err := h.service.GetForecasts(r.Context(), filters)
	if err != nil {
		log.Error("error getting forecasts", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, forecasts)
}

func (h *ForecastHandler) GetForecast(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid forecast ID", slog.String("error", err.Error()), slog.String("id", idStr))
		http.Error(w, "Invalid forecast ID", http.StatusBadRequest)
		return
	}

	log.Info("getting forecast by ID", slog.Int64("id", id))
	forecast, err := h.service.GetForecastByID(r.Context(), id)
	if err != nil {
		log.Error("error getting forecast by ID", slog.String("error", err.Error()), slog.Int64("id", id))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, forecast)
}

func (h *ForecastHandler) CreateForecast(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		log.Error("unauthorized", slog.String("error", "unauthorized"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var forecast models.Forecast
	if err := json.NewDecoder(r.Body).Decode(&forecast); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if forecast.Question == "" || forecast.ResolutionCriteria == "" || forecast.Category == "" {
		log.Error("question, resolution criteria, category, and closing date are required")
		http.Error(w, "Question, resolution criteria, category, and closing date are required", http.StatusBadRequest)
		return
	}

	// Set user ID from claims
	forecast.UserID = claims.UserID

	log.Info("creating forecast", slog.Any("forecast", forecast))
	err := h.service.CreateForecast(r.Context(), &forecast)
	if err != nil {
		log.Error("failed to create forecast", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, "forecast created")
}

func (h *ForecastHandler) DeleteForecast(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		ForecastID int64 `json:"forecast_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use UserID from claims
	if err := h.service.DeleteForecast(r.Context(), request.ForecastID, claims.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, "forecast deleted")
}

func (h *ForecastHandler) ResolveForecast(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		log.Error("unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var resolution struct {
		ID         int64  `json:"id"`
		Resolution string `json:"resolution"`
		Comment    string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resolution); err != nil {
		log.Error("invalid request body", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := claims.UserID

	log.Info("resolving forecast", slog.Any("resolution", resolution))
	if err := h.service.ResolveForecast(r.Context(),
		userID,
		resolution.ID,
		resolution.Resolution,
		resolution.Comment); err != nil {
		log.Error("failed to resolve forecast", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, "forecast resolved")
}

func (h *ForecastHandler) GetStaleAndNewForecasts(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Error("invalid user ID", slog.String("error", err.Error()))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Info("getting stale and new forecasts", slog.Int64("user_id", userID))
	forecasts, err := h.service.GetStaleAndNewForecasts(r.Context(), userID)
	if err != nil {
		log.Error("failed to get stale and new forecasts", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, forecasts)
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
