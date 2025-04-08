package handlers

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ForecastHandler struct {
	service *services.ForecastService
	cache   *cache.Cache
}

func NewForecastHandler(s *services.ForecastService, c *cache.Cache) *ForecastHandler {
	return &ForecastHandler{service: s, cache: c}
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

	category := "ALL"
	if request.Category != "" {
		category = request.Category
	}

	cacheKey := fmt.Sprintf("forecast:list:%s:%s", listType, category)

	if cachedList, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedList)
		return
	}

	forecasts, err := h.service.ForecastList(r.Context(), request.ListType, request.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.cache.Set(cacheKey, forecasts)

	respondJSON(w, http.StatusOK, forecasts)
}

func (h *ForecastHandler) GetForecast(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid forecast ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("forecast:detail:%d", id)

	if cachedForecast, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedForecast)
		return
	}

	forecast, err := h.service.GetForecastByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, forecast)

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

	h.cache.DeleteByPrefix("forecast:list:")
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

	h.cache.DeleteByPrefix("forecast:list:")
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

	// Use UserID from claims
	ownership, err := h.service.CheckForecastOwnership(r.Context(), resolution.ID, claims.UserID)
	if err == sql.ErrNoRows {
		http.Error(w, "forecast does not exist", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status, err := h.service.CheckForecastStatus(r.Context(), resolution.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !status {
		http.Error(w, "forecast is already resolved", http.StatusBadRequest)
		return
	}

	if !ownership {
		http.Error(w, "user does not own this forecast", http.StatusForbidden)
		return
	}

	if err := h.service.ResolveForecast(r.Context(),
		resolution.ID,
		resolution.Resolution,
		resolution.Comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	detailKey := fmt.Sprintf("forecast:detail:%d", resolution.ID)
	h.cache.Delete(detailKey)

	h.cache.DeleteByPrefix("forecast:list:")
	h.cache.DeleteByPrefix("scores:")

	respondJSON(w, http.StatusOK, "forecast resolved")
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
