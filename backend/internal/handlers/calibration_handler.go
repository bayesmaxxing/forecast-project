package handlers

import (
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/services"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CalibrationHandler struct {
	service *services.CalibrationService
}

func NewCalibrationHandler(s *services.CalibrationService) *CalibrationHandler {
	return &CalibrationHandler{service: s}
}

func (h *CalibrationHandler) GetCalibration(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	filters, err := parseCalibrationFilters(r)
	if err != nil {
		log.Error("invalid filter parameter", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("getting calibration data", slog.Any("filters", filters))
	data, err := h.service.GetCalibrationData(r.Context(), filters)
	if err != nil {
		log.Error("failed to get calibration data", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func (h *CalibrationHandler) GetCalibrationByUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	filters, err := parseCalibrationFilters(r)
	if err != nil {
		log.Error("invalid filter parameter", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("getting calibration data by users", slog.Any("filters", filters))
	data, err := h.service.GetCalibrationDataByUsers(r.Context(), filters)
	if err != nil {
		log.Error("failed to get calibration data by users", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, data)
}

func parseCalibrationFilters(r *http.Request) (models.CalibrationFilters, error) {
	queryParams := r.URL.Query()
	var filters models.CalibrationFilters

	userIDStr := queryParams.Get("user_id")
	if userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return filters, err
		}
		filters.UserID = &userID
	}

	categoryStr := queryParams.Get("category")
	if categoryStr != "" {
		category := strings.ToLower(categoryStr)
		filters.Category = &category
	}

	startDateStr := queryParams.Get("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return filters, err
		}
		filters.StartDate = &startDate
	}

	endDateStr := queryParams.Get("end_date")
	if endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return filters, err
		}
		filters.EndDate = &endDate
	}

	return filters, nil
}
