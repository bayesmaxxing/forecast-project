package handlers

import (
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

	cacheKey := fmt.Sprintf("forecasts_%s_%s", request.ListType, request.Category)

	if cachedList, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedList)
		return
	}
	forecasts, err := h.service.ForecastList(r.Context(), request.ListType, request.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, forecasts, 24*time.Hour)

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
	var forecast models.Forecast
	if err := json.NewDecoder(r.Body).Decode(&forecast); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateForecast(r.Context(), &forecast)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.DeleteByPrefix("forecasts")

	message := "forecast created"
	respondJSON(w, http.StatusCreated, message)
}

func (h *ForecastHandler) DeleteForecast(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID     int64 `json:"user_id"`
		ForecastID int64 `json:"forecast_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteForecast(r.Context(), request.ForecastID, request.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.DeleteByPrefix("forecasts")

	respondJSON(w, http.StatusOK, "forecast deleted")
}

func (h *ForecastHandler) ResolveForecast(w http.ResponseWriter, r *http.Request) {
	var resolution struct {
		ID         int64
		UserID     int64
		Resolution string
		Comment    string
	}

	if err := json.NewDecoder(r.Body).Decode(&resolution); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ownership, err := h.service.CheckForecastOwnership(r.Context(), resolution.ID, resolution.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	h.cache.DeleteByPrefix("forecasts")

	message := "forecast resolved"
	respondJSON(w, http.StatusOK, message)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
