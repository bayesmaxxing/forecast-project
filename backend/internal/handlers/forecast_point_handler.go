package handlers

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/services"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ForecastPointHandler struct {
	service *services.ForecastPointService
	cache   *cache.Cache
}

func NewForecastPointHandler(s *services.ForecastPointService, c *cache.Cache) *ForecastPointHandler {
	return &ForecastPointHandler{service: s, cache: c}
}

func (h *ForecastPointHandler) ListForecastPointsbyID(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/forecast-points/")
	forecastIDStr := strings.TrimSuffix(path, "/")

	if forecastIDStr == "" {
		http.Error(w, "Forecast ID is required", http.StatusBadRequest)
		return
	}

	forecastID, err := strconv.ParseInt(forecastIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid forecast ID", http.StatusBadRequest)
		return
	}

	points, err := h.service.GetForecastPointsByForecastID(r.Context(), forecastID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No forecast points found for this ID", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, points)
}

func (h *ForecastPointHandler) CreateForecastPoint(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var point models.ForecastPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the user ID from the JWT claims
	point.UserID = claims.UserID

	err := h.service.CreateForecastPoint(r.Context(), &point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.DeleteByPrefix("latest")

	message := "Forecast point created"

	respondJSON(w, http.StatusCreated, message)
}

func (h *ForecastPointHandler) ListAllForecastPoints(w http.ResponseWriter, r *http.Request) {
	points, err := h.service.GetAllForecastPoints(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, points)
}

func (h *ForecastPointHandler) ListLatestForecastPoints(w http.ResponseWriter, r *http.Request) {
	if cachedPoints, found := h.cache.Get("latest_all"); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	latestForecastPoints, err := h.service.GetLatestForecastPoints(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set("latest_all", latestForecastPoints, 24*time.Hour)

	respondJSON(w, http.StatusOK, latestForecastPoints)
}

func (h *ForecastPointHandler) ListLatestForecastPointsByUser(w http.ResponseWriter, r *http.Request) {
	var userID int64
	if err := json.NewDecoder(r.Body).Decode(&userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cachedPoints, found := h.cache.Get("latest_user_" + strconv.FormatInt(userID, 10)); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	latestForecastPoints, err := h.service.GetLatestForecastPointsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set("latest_user_"+strconv.FormatInt(userID, 10), latestForecastPoints, 24*time.Hour)

	respondJSON(w, http.StatusOK, latestForecastPoints)
}
