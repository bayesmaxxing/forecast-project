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

func (h *ForecastPointHandler) ListForecastPoints(w http.ResponseWriter, r *http.Request) {
	filters := models.PointFilters{}
	//parse query params
	queryParams := r.URL.Query()

	userIDstr := queryParams.Get("user_id")

	if userIDstr != "" {
		userID, err := strconv.ParseInt(userIDstr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user_id format", http.StatusBadRequest)
			return
		}
		filters.UserID = &userID
	}

	forecastIDstr := queryParams.Get("forecast_id")
	if forecastIDstr != "" {
		forecastID, err := strconv.ParseInt(forecastIDstr, 10, 64)
		if err != nil {
			http.Error(w, "invalid forecast_id format", http.StatusBadRequest)
			return
		}
		filters.ForecastID = &forecastID
	}

	dateStr := queryParams.Get("date")
	if dateStr != "" {
		date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
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

	points, err := h.service.GetForecastPoints(r.Context(), filters)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No forecast points found for these query parameters", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
