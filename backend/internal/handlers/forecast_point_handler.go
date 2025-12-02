package handlers

import (
	"backend/internal/auth"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/services"
	"database/sql"
	"encoding/json"
	"log/slog"
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

func (h *ForecastPointHandler) ListForecastPoints(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	filters := models.PointFilters{}
	//parse query params
	queryParams := r.URL.Query()

	userIDstr := queryParams.Get("user_id")

	if userIDstr != "" {
		userID, err := strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			log.Error("invalid user_id format", slog.String("error", err.Error()))
			http.Error(w, "invalid user_id format", http.StatusBadRequest)
			return
		}
		filters.UserID = &userID
	}

	forecastIDstr := queryParams.Get("forecast_id")
	if forecastIDstr != "" {
		forecastID, err := strconv.ParseInt(forecastIDstr, 10, 64)
		if err != nil {
			log.Error("invalid forecast_id format", slog.String("error", err.Error()))
			http.Error(w, "invalid forecast_id format", http.StatusBadRequest)
			return
		}
		filters.ForecastID = &forecastID
	}

	dateStr := queryParams.Get("date")
	if dateStr != "" {
		date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			log.Error("invalid date format", slog.String("error", err.Error()))
			http.Error(w, "invalid date format, expected YYYY-MM-DDTHH:MM:SSZ", http.StatusBadRequest)
			return
		}
		filters.Date = &date
	}

	distinctOnForecaststr := queryParams.Get("distinct")
	if distinctOnForecaststr != "" {
		distinctOnForecastBool, err := strconv.ParseBool(distinctOnForecaststr)
		if err != nil {
			http.Error(w, "invalid distinct format", http.StatusBadRequest)
			return
		}
		filters.DistinctOnForecast = &distinctOnForecastBool
	}

	orderByForecastIDstr := queryParams.Get("order_by_forecast_id")
	if orderByForecastIDstr != "" {
		orderByForecastIDBool, err := strconv.ParseBool(orderByForecastIDstr)
		if err != nil {
			http.Error(w, "invalid order_by_forecast_id format", http.StatusBadRequest)
			return
		}
		filters.OrderByForecastID = &orderByForecastIDBool
	}

	createdDirectionstr := queryParams.Get("created_direction")
	if createdDirectionstr != "" {
		createdDirectionString := strings.ToUpper(createdDirectionstr)
		if createdDirectionString != "ASC" && createdDirectionString != "DESC" {
			http.Error(w, "invalid created_direction format, expected ASC or DESC", http.StatusBadRequest)
			return
		} else if createdDirectionString == "ASC" {
			createdDirectionString = "ASC"
		} else {
			createdDirectionString = "DESC"
		}
		filters.CreatedDirection = &createdDirectionString
	}

	log.Info("getting forecast points", slog.Any("filters", filters))
	points, err := h.service.GetForecastPoints(r.Context(), filters)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("no forecast points found for these query parameters", slog.String("error", err.Error()))
			http.Error(w, "No forecast points found for these query parameters", http.StatusNotFound)
			return
		}
		log.Error("error getting forecast points", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, points)
}

func (h *ForecastPointHandler) CreateForecastPoint(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	// Get claims from context
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		log.Error("unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var point models.ForecastPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		log.Error("invalid request body", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the user ID from the JWT claims
	point.UserID = claims.UserID

	log.Info("creating forecast point", slog.Any("point", point))
	err := h.service.CreateForecastPoint(r.Context(), &point)
	if err != nil {
		log.Error("failed to create forecast point", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Forecast point created"

	respondJSON(w, http.StatusCreated, message)
}
