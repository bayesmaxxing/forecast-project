package handlers_test

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/handlers"
	"backend/internal/models"
	"backend/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// MockForecastService implements services.ForecastService interface for testing
type MockForecastService struct {
	forecasts               map[int64]*models.Forecast
	nextID                  int64
	forecastListErr         error
	forecastListResult      []*models.Forecast
	getForecastByIDErr      error
	createForecastErr       error
	deleteForecastErr       error
	resolveForecastErr      error
	checkForecastOwnership  bool
	checkForecastOwnershipErr error
}

func NewMockForecastService() *services.ForecastService {
	mock := &MockForecastService{
		forecasts: make(map[int64]*models.Forecast),
		nextID:    1,
	}
	var service services.ForecastService = mock
	return &service
}

func (m *MockForecastService) ForecastList(ctx context.Context, listType string, category string) ([]*models.Forecast, error) {
	if m.forecastListErr != nil {
		return nil, m.forecastListErr
	}
	if m.forecastListResult != nil {
		return m.forecastListResult, nil
	}
	var result []*models.Forecast
	for _, f := range m.forecasts {
		if (listType == "open" && !f.IsResolved()) || (listType == "resolved" && f.IsResolved()) {
			if category == "" || f.Category == category {
				result = append(result, f)
			}
		}
	}
	return result, nil
}

func (m *MockForecastService) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
	if m.getForecastByIDErr != nil {
		return nil, m.getForecastByIDErr
	}
	f, exists := m.forecasts[id]
	if !exists {
		return nil, errors.New("forecast not found")
	}
	return f, nil
}

func (m *MockForecastService) CreateForecast(ctx context.Context, f *models.Forecast) error {
	if m.createForecastErr != nil {
		return m.createForecastErr
	}
	f.ID = m.nextID
	m.forecasts[f.ID] = f
	m.nextID++
	return nil
}

func (m *MockForecastService) DeleteForecast(ctx context.Context, id int64, userID int64) error {
	if m.deleteForecastErr != nil {
		return m.deleteForecastErr
	}
	_, exists := m.forecasts[id]
	if !exists {
		return errors.New("forecast not found")
	}
	delete(m.forecasts, id)
	return nil
}

func (m *MockForecastService) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	_, exists := m.forecasts[f.ID]
	if !exists {
		return errors.New("forecast not found")
	}
	m.forecasts[f.ID] = f
	return nil
}

func (m *MockForecastService) ResolveForecast(ctx context.Context, id int64, resolution string, comment string) error {
	if m.resolveForecastErr != nil {
		return m.resolveForecastErr
	}
	f, exists := m.forecasts[id]
	if !exists {
		return errors.New("forecast not found")
	}
	now := time.Now()
	f.ResolvedAt = &now
	f.Resolution = &resolution
	f.ResolutionComment = &comment
	return nil
}

func (m *MockForecastService) CheckForecastOwnership(ctx context.Context, id int64, userID int64) (bool, error) {
	if m.checkForecastOwnershipErr != nil {
		return false, m.checkForecastOwnershipErr
	}
	return m.checkForecastOwnership, nil
}

// Helper function to get the mock service from the interface
func getMockForecastService(s *services.ForecastService) *MockForecastService {
	return (*s).(*MockForecastService)
}

// MockCache implements the cache.Cache interface for testing
type MockCache struct {
	data map[string]interface{}
}

func NewMockCache() *cache.Cache {
	mock := &MockCache{
		data: make(map[string]interface{}),
	}
	var c cache.Cache = mock
	return &c
}

func (c *MockCache) Get(key string) (interface{}, bool) {
	value, found := c.data[key]
	return value, found
}

func (c *MockCache) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *MockCache) Delete(key string) {
	delete(c.data, key)
}

func (c *MockCache) DeleteByPrefix(prefix string) {
	for key := range c.data {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.data, key)
		}
	}
}

// Helper function to get the mock cache from the interface
func getMockCache(c *cache.Cache) *MockCache {
	return (*c).(*MockCache)
}

// Helper function to create a request with JSON body
func createRequestWithJSON(method, url string, body interface{}) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Helper function to add auth context to a request
func addAuthContext(req *http.Request, userID int64) *http.Request {
	claims := &auth.Claims{
		UserID: userID,
	}
	ctx := context.WithValue(req.Context(), auth.UserContextKey, claims)
	return req.WithContext(ctx)
}

// Test for ListForecasts handler
func TestListForecasts(t *testing.T) {
	mockService := NewMockForecastService()
	mockCache := NewMockCache()
	handler := handlers.NewForecastHandler(mockService, mockCache)

	// Create test data
	forecast1 := &models.Forecast{
		ID:       1,
		Question: "Test Question 1",
		Category: "Category1",
	}
	forecast2 := &models.Forecast{
		ID:       2,
		Question: "Test Question 2",
		Category: "Category2",
	}
	mockService.forecasts[forecast1.ID] = forecast1
	mockService.forecasts[forecast2.ID] = forecast2

	// Setup mock return for forecast list
	mockService.forecastListResult = []*models.Forecast{forecast1, forecast2}

	t.Run("SuccessfulList", func(t *testing.T) {
		// Create request with body
		body := map[string]string{
			"list_type": "open",
			"category":  "",
		}
		req := createRequestWithJSON(http.MethodPost, "/forecasts", body)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ListForecasts(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify response body contains forecasts
		var response []*models.Forecast
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		if len(response) != 2 {
			t.Errorf("Expected 2 forecasts, got %d", len(response))
		}
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		// Create request with invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/forecasts", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ListForecasts(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Setup mock to return error
		mockService.forecastListErr = errors.New("service error")

		// Create request with body
		body := map[string]string{
			"list_type": "open",
			"category":  "",
		}
		req := createRequestWithJSON(http.MethodPost, "/forecasts", body)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ListForecasts(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Reset mock
		mockService.forecastListErr = nil
	})

	t.Run("CachedResponse", func(t *testing.T) {
		// Set up cache with forecasts
		cacheKey := "forecasts_open_"
		cachedForecasts := []*models.Forecast{forecast1}
		mockCache.Set(cacheKey, cachedForecasts)

		// Create request with body
		body := map[string]string{
			"list_type": "open",
			"category":  "",
		}
		req := createRequestWithJSON(http.MethodPost, "/forecasts", body)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ListForecasts(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify response body contains cached forecasts
		var response []*models.Forecast
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		if len(response) != 1 {
			t.Errorf("Expected 1 forecast from cache, got %d", len(response))
		}
	})
}

// Test for GetForecast handler
func TestGetForecast(t *testing.T) {
	mockService := NewMockForecastService()
	mockCache := NewMockCache()
	handler := handlers.NewForecastHandler(mockService, mockCache)

	// Create test data
	forecast := &models.Forecast{
		ID:       1,
		Question: "Test Question",
		Category: "Test Category",
	}
	mockService.forecasts[forecast.ID] = forecast

	t.Run("SuccessfulGet", func(t *testing.T) {
		// Create test request
		req, _ := http.NewRequest(http.MethodGet, "/forecasts/"+strconv.FormatInt(forecast.ID, 10), nil)
		rr := httptest.NewRecorder()

		// Create a router to capture pathValue
		req = req.WithPathValue("id", strconv.FormatInt(forecast.ID, 10))

		// Call the handler
		handler.GetForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify response body contains forecast
		var response models.Forecast
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error unmarshaling response: %v", err)
		}
		if response.ID != forecast.ID {
			t.Errorf("Expected forecast ID %d, got %d", forecast.ID, response.ID)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		// Create test request with invalid ID
		req, _ := http.NewRequest(http.MethodGet, "/forecasts/invalid", nil)
		rr := httptest.NewRecorder()

		// Create a router to capture pathValue
		req = req.WithPathValue("id", "invalid")

		// Call the handler
		handler.GetForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("ForecastNotFound", func(t *testing.T) {
		// Setup mock to return error
		mockService.getForecastByIDErr = errors.New("forecast not found")

		// Create test request
		req, _ := http.NewRequest(http.MethodGet, "/forecasts/999", nil)
		rr := httptest.NewRecorder()

		// Create a router to capture pathValue
		req = req.WithPathValue("id", "999")

		// Call the handler
		handler.GetForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Reset mock
		mockService.getForecastByIDErr = nil
	})
}

// Test for CreateForecast handler
func TestCreateForecast(t *testing.T) {
	mockService := NewMockForecastService()
	mockCache := NewMockCache()
	handler := handlers.NewForecastHandler(mockService, mockCache)

	t.Run("SuccessfulCreate", func(t *testing.T) {
		// Create forecast data
		forecast := models.Forecast{
			Question: "New Test Question",
			Category: "Test Category",
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodPost, "/forecasts", forecast)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CreateForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		// Verify cache was cleared
		_, found := mockCache.Get("forecasts_open_")
		if found {
			t.Error("Cache was not cleared after creating forecast")
		}
	})

	t.Run("NoAuthentication", func(t *testing.T) {
		// Create request without auth context
		forecast := models.Forecast{
			Question: "New Test Question",
			Category: "Test Category",
		}
		req := createRequestWithJSON(http.MethodPost, "/forecasts", forecast)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CreateForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		// Create request with invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/forecasts", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CreateForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Setup mock to return error
		mockService.createForecastErr = errors.New("service error")

		// Create forecast data
		forecast := models.Forecast{
			Question: "New Test Question",
			Category: "Test Category",
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodPost, "/forecasts", forecast)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CreateForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Reset mock
		mockService.createForecastErr = nil
	})
}

// Test for DeleteForecast handler
func TestDeleteForecast(t *testing.T) {
	mockService := NewMockForecastService()
	mockCache := NewMockCache()
	handler := handlers.NewForecastHandler(mockService, mockCache)

	// Create test data
	forecast := &models.Forecast{
		ID:       1,
		Question: "Test Question",
		Category: "Test Category",
		UserID:   1,
	}
	mockService.forecasts[forecast.ID] = forecast

	t.Run("SuccessfulDelete", func(t *testing.T) {
		// Create delete request data
		deleteRequest := struct {
			ForecastID int64 `json:"forecast_id"`
		}{
			ForecastID: forecast.ID,
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodDelete, "/forecasts", deleteRequest)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.DeleteForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify cache was cleared
		_, found := mockCache.Get("forecasts_open_")
		if found {
			t.Error("Cache was not cleared after deleting forecast")
		}
	})

	t.Run("NoAuthentication", func(t *testing.T) {
		// Create request without auth context
		deleteRequest := struct {
			ForecastID int64 `json:"forecast_id"`
		}{
			ForecastID: forecast.ID,
		}
		req := createRequestWithJSON(http.MethodDelete, "/forecasts", deleteRequest)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.DeleteForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		// Create request with invalid JSON
		req, _ := http.NewRequest(http.MethodDelete, "/forecasts", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.DeleteForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Setup mock to return error
		mockService.deleteForecastErr = errors.New("service error")

		// Create delete request data
		deleteRequest := struct {
			ForecastID int64 `json:"forecast_id"`
		}{
			ForecastID: forecast.ID,
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodDelete, "/forecasts", deleteRequest)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.DeleteForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Reset mock
		mockService.deleteForecastErr = nil
	})
}

// Test for ResolveForecast handler
func TestResolveForecast(t *testing.T) {
	mockService := NewMockForecastService()
	mockCache := NewMockCache()
	handler := handlers.NewForecastHandler(mockService, mockCache)

	// Create test data
	forecast := &models.Forecast{
		ID:       1,
		Question: "Test Question",
		Category: "Test Category",
		UserID:   1,
	}
	mockService.forecasts[forecast.ID] = forecast

	t.Run("SuccessfulResolve", func(t *testing.T) {
		// Set ownership check to return true
		mockService.checkForecastOwnership = true

		// Create resolve request data
		resolveRequest := struct {
			ID         int64  `json:"id"`
			Resolution string `json:"resolution"`
			Comment    string `json:"comment"`
		}{
			ID:         forecast.ID,
			Resolution: "1",
			Comment:    "Test resolution comment",
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodPut, "/forecasts/resolve", resolveRequest)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ResolveForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify cache was cleared
		_, found := mockCache.Get("forecasts_open_")
		if found {
			t.Error("Cache was not cleared after resolving forecast")
		}
	})

	t.Run("NoAuthentication", func(t *testing.T) {
		// Create request without auth context
		resolveRequest := struct {
			ID         int64  `json:"id"`
			Resolution string `json:"resolution"`
			Comment    string `json:"comment"`
		}{
			ID:         forecast.ID,
			Resolution: "1",
			Comment:    "Test resolution comment",
		}
		req := createRequestWithJSON(http.MethodPut, "/forecasts/resolve", resolveRequest)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ResolveForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		// Create request with invalid JSON
		req, _ := http.NewRequest(http.MethodPut, "/forecasts/resolve", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ResolveForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("NotOwner", func(t *testing.T) {
		// Set ownership check to return false
		mockService.checkForecastOwnership = false

		// Create resolve request data
		resolveRequest := struct {
			ID         int64  `json:"id"`
			Resolution string `json:"resolution"`
			Comment    string `json:"comment"`
		}{
			ID:         forecast.ID,
			Resolution: "1",
			Comment:    "Test resolution comment",
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodPut, "/forecasts/resolve", resolveRequest)
		req = addAuthContext(req, 2) // Different user ID
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ResolveForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
		}
	})

	t.Run("OwnershipCheckError", func(t *testing.T) {
		// Setup mock to return error for ownership check
		mockService.checkForecastOwnershipErr = errors.New("database error")

		// Create resolve request data
		resolveRequest := struct {
			ID         int64  `json:"id"`
			Resolution string `json:"resolution"`
			Comment    string `json:"comment"`
		}{
			ID:         forecast.ID,
			Resolution: "1",
			Comment:    "Test resolution comment",
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodPut, "/forecasts/resolve", resolveRequest)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ResolveForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Reset mock
		mockService.checkForecastOwnershipErr = nil
	})

	t.Run("ResolveError", func(t *testing.T) {
		// Set ownership check to return true
		mockService.checkForecastOwnership = true
		
		// Setup mock to return error
		mockService.resolveForecastErr = errors.New("service error")

		// Create resolve request data
		resolveRequest := struct {
			ID         int64  `json:"id"`
			Resolution string `json:"resolution"`
			Comment    string `json:"comment"`
		}{
			ID:         forecast.ID,
			Resolution: "1",
			Comment:    "Test resolution comment",
		}

		// Create request with auth context
		req := createRequestWithJSON(http.MethodPut, "/forecasts/resolve", resolveRequest)
		req = addAuthContext(req, 1)
		rr := httptest.NewRecorder()

		// Call the handler
		handler.ResolveForecast(rr, req)

		// Check response
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Reset mock
		mockService.resolveForecastErr = nil
	})
}

// Test the respondJSON function (indirectly)
func TestRespondJSON(t *testing.T) {
	mockService := NewMockForecastService()
	mockCache := NewMockCache()
	handler := handlers.NewForecastHandler(mockService, mockCache)

	// Create test data
	forecast := &models.Forecast{
		ID:       1,
		Question: "Test Question",
		Category: "Test Category",
		UserID:   1,
	}
	mockService.forecasts[forecast.ID] = forecast

	// Create test request
	req, _ := http.NewRequest(http.MethodGet, "/forecasts/"+strconv.FormatInt(forecast.ID, 10), nil)
	rr := httptest.NewRecorder()

	// Create a router to capture pathValue
	req = req.WithPathValue("id", strconv.FormatInt(forecast.ID, 10))

	// Call the handler
	handler.GetForecast(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify content type is application/json
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}
}