package handlers

import (
	"backend/internal/auth"
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
}

func NewForecastPointHandler(s *services.ForecastPointService) *ForecastPointHandler {
	return &ForecastPointHandler{service: s}
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

func (h *ForecastPointHandler) ListForecastPointsbyIDAndUser(w http.ResponseWriter, r *http.Request) {
	var point_request struct {
		ForecastID int64 `json:"forecast_id"`
		UserID     int64 `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&point_request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	latestForecastPoints, err := h.service.GetLatestForecastPoints(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	latestForecastPoints, err := h.service.GetLatestForecastPointsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	userIDStr := r.URL.Query().Get("user_id")
	userId, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	points, err := h.service.GetOrderedForecastPoints(r.Context(), forecastID, userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, points)
}

func (h *ForecastPointHandler) ListForecastPointsByDate(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from path
	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id path parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

	// Parse optional date query parameter (format: YYYY-MM-DD)
	var date *time.Time
	dateStr := r.URL.Query().Get("date")
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		date = &parsedDate
	}

	points, err := h.service.GetForecastPointsByDate(r.Context(), userID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, points)
}
