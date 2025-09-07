package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
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
	var request struct {
		ListType string `json:"list_type"`
		Category string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle empty values with defaults
	listType := "ALL"
	if request.ListType != "" {
		listType = request.ListType
	}

	// Validate list type
	validListTypes := map[string]bool{"open": true, "resolved": true, "ALL": true}
	if !validListTypes[listType] {
		http.Error(w, "invalid list_type", http.StatusBadRequest)
		return
	}

	category := "ALL"
	if request.Category != "" {
		category = request.Category
	}

	forecasts, err := h.service.ForecastList(r.Context(), listType, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusOK, forecasts)
}

func (h *ForecastHandler) GetForecast(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid forecast ID", http.StatusBadRequest)
		return
	}

	forecast, err := h.service.GetForecastByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, forecast)
}

func (h *ForecastHandler) CreateForecast(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var forecast models.Forecast
	if err := json.NewDecoder(r.Body).Decode(&forecast); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if forecast.Question == "" || forecast.ResolutionCriteria == "" || forecast.Category == "" {
		http.Error(w, "Question, resolution criteria, and category are required", http.StatusBadRequest)
		return
	}

	// Set user ID from claims
	forecast.UserID = claims.UserID

	err := h.service.CreateForecast(r.Context(), &forecast)
	if err != nil {
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
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var resolution struct {
		ID         int64  `json:"id"`
		Resolution string `json:"resolution"`
		Comment    string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resolution); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := claims.UserID

	if err := h.service.ResolveForecast(r.Context(),
		userID,
		resolution.ID,
		resolution.Resolution,
		resolution.Comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, "forecast resolved")
}

func (h *ForecastHandler) GetStaleAndNewForecasts(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	forecasts, err := h.service.GetStaleAndNewForecasts(r.Context(), userID)
	if err != nil {
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
