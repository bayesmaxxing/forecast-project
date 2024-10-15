package handlers

import (
	"encoding/json"
	"go_api/internal/models"
	"go_api/internal/services"
	"net/http"
	"strconv"
)

type ForecastHandler struct {
	service *services.ForecastService
}

func NewForecastHandler(s *services.ForecastService) *ForecastHandler {
	return &ForecastHandler{service: s}
}

func (h *ForecastHandler) ListForecasts(w http.ResponseWriter, r *http.Request) {
	listType := r.URL.Query().Get("type")
	category := r.URL.Query().Get("category")

	forecasts, err := h.service.ForecastList(r.Context(), listType, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	message := "forecast created"
	respondJSON(w, http.StatusCreated, message)
}

func (h *ForecastHandler) DeleteForecast(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid forecast ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteForecast(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ForecastHandler) ResolveForecast(w http.ResponseWriter, r *http.Request) {
	var resolution struct {
		ID         int64
		Resolution string
		Comment    string
	}

	if err := json.NewDecoder(r.Body).Decode(&resolution); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ResolveForecast(r.Context(),
		resolution.ID,
		resolution.Resolution,
		resolution.Comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "forecast resolved"
	respondJSON(w, http.StatusOK, message)
}

func (h *ForecastHandler) GetAggregatedScores(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	scores, err := h.service.GetAggregatedScores(r.Context(), category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, scores)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
