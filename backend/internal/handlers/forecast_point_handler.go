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
	"strings"
	"time"
)

type ForecastPointHandler struct {
	service         *services.ForecastPointService
	forecastService *services.ForecastService
	cache           *cache.Cache
}

func NewForecastPointHandler(s *services.ForecastPointService, f *services.ForecastService, c *cache.Cache) *ForecastPointHandler {
	return &ForecastPointHandler{service: s, forecastService: f, cache: c}
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

	cacheKey := fmt.Sprintf("point:list:%d", forecastID)
	if cachedPoints, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedPoints)
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

	h.cache.Set(cacheKey, points)

	respondJSON(w, http.StatusOK, points)
}

func (h *ForecastPointHandler) ListForecastPointsbyIDAndUser(w http.ResponseWriter, r *http.Request) {
	var point_request struct {
		ForecastID int64 `json:"forecast_id"`
		UserID     int64 `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&point_request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("point:list:user:%d:%d", point_request.UserID, point_request.ForecastID)

	if cachedPoints, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	points, err := h.service.GetForecastPointsByForecastIDAndUser(r.Context(), point_request.ForecastID, point_request.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No forecast points found for this ID", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, points)

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

	// Check if the forecast exists
	_, err := h.forecastService.GetForecastByID(r.Context(), point.ForecastID)
	if err != nil {
		http.Error(w, "Forecast not found", http.StatusBadRequest)
		return
	}

	err = h.service.CreateForecastPoint(r.Context(), &point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Delete(fmt.Sprintf("point:list:%d", point.ForecastID))
	h.cache.Delete("point:all:latest")
	h.cache.Delete("point:all")
	message := "Forecast point created"

	respondJSON(w, http.StatusCreated, message)
}

func (h *ForecastPointHandler) ListAllForecastPoints(w http.ResponseWriter, r *http.Request) {
	cacheKey := "point:all"
	if cachedPoints, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	points, err := h.service.GetAllForecastPoints(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, points)

	respondJSON(w, http.StatusOK, points)
}

func (h *ForecastPointHandler) ListLatestForecastPoints(w http.ResponseWriter, r *http.Request) {
	cacheKey := "point:all:latest"
	if cachedPoints, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	latestForecastPoints, err := h.service.GetLatestForecastPoints(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, latestForecastPoints)

	respondJSON(w, http.StatusOK, latestForecastPoints)
}

func (h *ForecastPointHandler) ListLatestForecastPointsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("point:all:latest:%d", userID)
	if cachedPoints, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	latestForecastPoints, err := h.service.GetLatestForecastPointsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.cache.Set(cacheKey, latestForecastPoints)

	respondJSON(w, http.StatusOK, latestForecastPoints)
}

type GraphPoint struct {
	PointForecast float64   `json:"point_forecast"`
	CreatedAt     time.Time `json:"created"`
	UserID        int64     `json:"user_id"`
}

func (h *ForecastPointHandler) ListOrderedForecastPoints(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/forecast-points/ordered/")
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

	// Check if user_id query parameter is provided
	userIDStr := r.URL.Query().Get("user_id")

	userId := "ALL"
	if userIDStr != "" {
		userId = userIDStr
	}

	cacheKey := fmt.Sprintf("point:ordered:%s:%d", userId, forecastID)
	if cachedPoints, found := h.cache.Get(cacheKey); found {
		respondJSON(w, http.StatusOK, cachedPoints)
		return
	}

	var points []*models.ForecastPoint

	if userIDStr != "" {
		// If user_id is provided, return points for that user only
		var userID int64
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		points, err = h.service.GetForecastPointsByForecastIDAndUser(r.Context(), forecastID, userID)
	} else {
		// If no user_id is provided, return all points
		points, err = h.service.GetOrderedForecastPointsByForecastID(r.Context(), forecastID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No forecast points found for this ID", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Transform the points to a simpler format for graphing
	graphPoints := make([]GraphPoint, len(points))
	for i, p := range points {
		graphPoints[i] = GraphPoint{
			PointForecast: p.PointForecast,
			CreatedAt:     p.CreatedAt,
			UserID:        p.UserID,
		}
	}

	h.cache.Set(cacheKey, graphPoints)

	respondJSON(w, http.StatusOK, graphPoints)
}
