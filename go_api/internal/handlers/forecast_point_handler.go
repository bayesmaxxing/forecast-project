package handlers

import (
	"database/sql"
	"encoding/json"
	"go_api/internal/models"
	"go_api/internal/services"
	"net/http"
	"strconv"
	"strings"
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

func (h *ForecastPointHandler) CreateForecastPoint(w http.ResponseWriter, r *http.Request) {
	var point models.ForecastPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
