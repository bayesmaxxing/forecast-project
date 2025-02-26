package services_test

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"errors"
	"testing"
	"time"
)

// ForecastRepositoryMock is a mock implementation of the repository.ForecastRepository interface
type ForecastRepositoryMock struct {
	forecasts       map[int64]*models.Forecast
	nextID          int64
	getForecastErr  error
	createErr       error
	updateErr       error
	deleteErr       error
	ownershipErr    error
	isOwner         bool
	listOpenErr     error
	listResolvedErr error
}

func NewForecastRepositoryMock() repository.ForecastRepository {
	return &ForecastRepositoryMock{
		forecasts: make(map[int64]*models.Forecast),
		nextID:    1,
	}
}

func (m *ForecastRepositoryMock) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
	if m.getForecastErr != nil {
		return nil, m.getForecastErr
	}
	f, exists := m.forecasts[id]
	if !exists {
		return nil, errors.New("forecast not found")
	}
	return f, nil
}

func (m *ForecastRepositoryMock) CheckForecastOwnership(ctx context.Context, id int64, user_id int64) (bool, error) {
	if m.ownershipErr != nil {
		return false, m.ownershipErr
	}
	return m.isOwner, nil
}

func (m *ForecastRepositoryMock) CreateForecast(ctx context.Context, f *models.Forecast) error {
	if m.createErr != nil {
		return m.createErr
	}
	f.ID = m.nextID
	m.nextID++
	m.forecasts[f.ID] = f
	return nil
}

func (m *ForecastRepositoryMock) DeleteForecast(ctx context.Context, id int64, user_id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.forecasts, id)
	return nil
}

func (m *ForecastRepositoryMock) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.forecasts[f.ID] = f
	return nil
}

func (m *ForecastRepositoryMock) ListOpenForecasts(ctx context.Context) ([]*models.Forecast, error) {
	if m.listOpenErr != nil {
		return nil, m.listOpenErr
	}
	var result []*models.Forecast
	for _, f := range m.forecasts {
		if !f.IsResolved() {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *ForecastRepositoryMock) ListResolvedForecasts(ctx context.Context) ([]*models.Forecast, error) {
	if m.listResolvedErr != nil {
		return nil, m.listResolvedErr
	}
	var result []*models.Forecast
	for _, f := range m.forecasts {
		if f.IsResolved() {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *ForecastRepositoryMock) ListOpenForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	if m.listOpenErr != nil {
		return nil, m.listOpenErr
	}

	return m.ListOpenForecasts(ctx)
}

func (m *ForecastRepositoryMock) ListResolvedForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	if m.listResolvedErr != nil {
		return nil, m.listResolvedErr
	}

	return m.ListResolvedForecasts(ctx)
}

// Helper function to access the underlying mock
func getForecastRepoMock(repo repository.ForecastRepository) *ForecastRepositoryMock {
	return repo.(*ForecastRepositoryMock)
}

// Test for GetForecastByID
func TestGetForecastByID(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Create a test forecast
	testForecast := &models.Forecast{
		Question: "Test Question",
		Category: "Test Category",
	}
	getForecastRepoMock(mockRepo).CreateForecast(ctx, testForecast)

	// Test successful retrieval
	forecast, err := service.GetForecastByID(ctx, testForecast.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if forecast.ID != testForecast.ID {
		t.Errorf("Expected forecast ID %d, got %d", testForecast.ID, forecast.ID)
	}

	// Test error case
	getForecastRepoMock(mockRepo).getForecastErr = errors.New("database error")
	_, err = service.GetForecastByID(ctx, testForecast.ID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for CheckForecastOwnership
func TestCheckForecastOwnership(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Test when user is the owner
	getForecastRepoMock(mockRepo).isOwner = true
	isOwner, err := service.CheckForecastOwnership(ctx, 1, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !isOwner {
		t.Error("Expected true for ownership check, got false")
	}

	// Test when user is not the owner
	getForecastRepoMock(mockRepo).isOwner = false
	isOwner, err = service.CheckForecastOwnership(ctx, 1, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isOwner {
		t.Error("Expected false for ownership check, got true")
	}

	// Test error case
	getForecastRepoMock(mockRepo).ownershipErr = errors.New("database error")
	_, err = service.CheckForecastOwnership(ctx, 1, 1)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for CreateForecast
func TestCreateForecast(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Test successful creation
	testForecast := &models.Forecast{
		Question: "Test Question",
		Category: "Test Category",
	}
	err := service.CreateForecast(ctx, testForecast)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if testForecast.ID != 1 {
		t.Errorf("Expected forecast ID 1, got %d", testForecast.ID)
	}

	// Test error case
	getForecastRepoMock(mockRepo).createErr = errors.New("database error")
	err = service.CreateForecast(ctx, testForecast)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for DeleteForecast
func TestDeleteForecast(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Create a test forecast
	testForecast := &models.Forecast{
		Question: "Test Question",
		Category: "Test Category",
	}
	getForecastRepoMock(mockRepo).CreateForecast(ctx, testForecast)

	// Test successful deletion
	err := service.DeleteForecast(ctx, testForecast.ID, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test error case
	getForecastRepoMock(mockRepo).deleteErr = errors.New("database error")
	err = service.DeleteForecast(ctx, testForecast.ID, 1)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for UpdateForecast
func TestUpdateForecast(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Create a test forecast
	testForecast := &models.Forecast{
		Question: "Test Question",
		Category: "Test Category",
	}
	getForecastRepoMock(mockRepo).CreateForecast(ctx, testForecast)

	// Update the forecast
	testForecast.Question = "Updated Question"
	err := service.UpdateForecast(ctx, testForecast)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the update
	updatedForecast, _ := service.GetForecastByID(ctx, testForecast.ID)
	if updatedForecast.Question != "Updated Question" {
		t.Errorf("Expected question 'Updated Question', got '%s'", updatedForecast.Question)
	}

	// Test error case
	getForecastRepoMock(mockRepo).updateErr = errors.New("database error")
	err = service.UpdateForecast(ctx, testForecast)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for ResolveForecast
func TestResolveForecast(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Create a test forecast
	testForecast := &models.Forecast{
		ID:       1,
		Question: "Test Question",
		Category: "Test Category",
	}
	getForecastRepoMock(mockRepo).forecasts[testForecast.ID] = testForecast

	t.Run("ResolveForecastWithTrueOutcome", func(t *testing.T) {
		// Reset the point repository and add points from two users
		getMockForecastPointRepo(mockPointRepo).points = []*models.ForecastPoint{
			{
				ForecastID:    1,
				PointForecast: 0.7,
				UserID:        1,
				CreatedAt:     time.Now(),
			},
			{
				ForecastID:    1,
				PointForecast: 0.8,
				UserID:        2,
				CreatedAt:     time.Now(),
			},
		}

		// Reset scores repository
		getMockScoreRepository(mockScoreRepo).scores = []models.Scores{}

		// Test resolution with outcome "1" (true)
		resolution := "1"
		comment := "Test resolution"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, comment)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify forecast is resolved
		resolvedForecast, _ := service.GetForecastByID(ctx, testForecast.ID)
		if !resolvedForecast.IsResolved() {
			t.Error("Expected forecast to be resolved, but it's not")
		}
		if *resolvedForecast.Resolution != resolution {
			t.Errorf("Expected resolution '%s', got '%s'", resolution, *resolvedForecast.Resolution)
		}
		if *resolvedForecast.ResolutionComment != comment {
			t.Errorf("Expected comment '%s', got '%s'", comment, *resolvedForecast.ResolutionComment)
		}

		// Verify scores were created for both users
		scores := getMockScoreRepository(mockScoreRepo).scores
		if len(scores) != 2 {
			t.Errorf("Expected 2 scores to be created, got %d", len(scores))
		}

		// Check if scores are for the right users and forecast
		userIDs := make(map[int64]bool)
		for _, score := range scores {
			userIDs[score.UserID] = true
			if score.ForecastID != testForecast.ID {
				t.Errorf("Expected score for forecast %d, got %d", testForecast.ID, score.ForecastID)
			}
		}
		if len(userIDs) != 2 || !userIDs[1] || !userIDs[2] {
			t.Error("Expected scores for users 1 and 2")
		}
	})

	t.Run("ResolveForecastWithFalseOutcome", func(t *testing.T) {
		// Reset the state
		getMockScoreRepository(mockScoreRepo).scores = []models.Scores{}

		// Test resolution with outcome "0" (false)
		resolution := "0"
		comment := "Test resolution false"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, comment)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify scores were created
		scores := getMockScoreRepository(mockScoreRepo).scores
		if len(scores) != 2 {
			t.Errorf("Expected 2 scores to be created, got %d", len(scores))
		}
	})

	t.Run("ResolveForecastWithDash", func(t *testing.T) {
		// Reset the state
		getMockScoreRepository(mockScoreRepo).scores = []models.Scores{}

		// Test resolution with outcome "-" (no scoring)
		resolution := "-"
		comment := "Test resolution dash"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, comment)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify NO scores were created
		scores := getMockScoreRepository(mockScoreRepo).scores
		if len(scores) != 0 {
			t.Errorf("Expected 0 scores to be created, got %d", len(scores))
		}
	})

	t.Run("NoForecastPoints", func(t *testing.T) {
		// Test no forecast points case
		getMockForecastPointRepo(mockPointRepo).points = []*models.ForecastPoint{}
		resolution := "1"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, "Test comment")
		if err == nil || err.Error() != "no forecast points found" {
			t.Errorf("Expected 'no forecast points found' error, got %v", err)
		}
	})

	t.Run("GetForecastError", func(t *testing.T) {
		// Test error from GetForecastByID
		getForecastRepoMock(mockRepo).getForecastErr = errors.New("database error")
		resolution := "1"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, "Test comment")
		if err == nil {
			t.Error("Expected error from GetForecastByID, got nil")
		}
		// Reset for subsequent tests
		getForecastRepoMock(mockRepo).getForecastErr = nil
	})

	t.Run("GetPointsError", func(t *testing.T) {
		// Test error from GetForecastPointsByForecastID
		getMockForecastPointRepo(mockPointRepo).getByForecastIDErr = errors.New("database error")
		resolution := "1"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, "Test comment")
		if err == nil {
			t.Error("Expected error from GetForecastPointsByForecastID, got nil")
		}
		// Reset for subsequent tests
		getMockForecastPointRepo(mockPointRepo).getByForecastIDErr = nil
	})

	t.Run("CreateScoreError", func(t *testing.T) {
		// Reset points
		getMockForecastPointRepo(mockPointRepo).points = []*models.ForecastPoint{
			{
				ForecastID:    1,
				PointForecast: 0.7,
				UserID:        1,
				CreatedAt:     time.Now(),
			},
		}

		// Test error from CreateScore
		getMockScoreRepository(mockScoreRepo).createScoreErr = errors.New("database error")
		resolution := "1"
		err := service.ResolveForecast(ctx, testForecast.ID, resolution, "Test comment")
		if err == nil {
			t.Error("Expected error from CreateScore, got nil")
		}
		// Reset for subsequent tests
		getMockScoreRepository(mockScoreRepo).createScoreErr = nil
	})
}

// Test for ForecastList
func TestForecastList(t *testing.T) {
	// Setup
	mockRepo := NewForecastRepositoryMock()
	mockPointRepo := NewMockForecastPointRepository()
	mockScoreRepo := NewMockScoreRepository()
	service := services.NewForecastService(mockRepo, mockPointRepo, mockScoreRepo)
	ctx := context.Background()

	// Create test forecasts
	openForecast := &models.Forecast{
		Question: "Open Question",
		Category: "Test Category",
	}
	getForecastRepoMock(mockRepo).CreateForecast(ctx, openForecast)

	resolvedForecast := &models.Forecast{
		Question: "Resolved Question",
		Category: "Test Category",
	}
	getForecastRepoMock(mockRepo).CreateForecast(ctx, resolvedForecast)

	// Mark the second forecast as resolved
	now := time.Now()
	resolution := "1"
	resolvedForecast.ResolvedAt = &now
	resolvedForecast.Resolution = &resolution
	getForecastRepoMock(mockRepo).UpdateForecast(ctx, resolvedForecast)

	// Test listing open forecasts
	openForecasts, err := service.ForecastList(ctx, "open", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(openForecasts) != 1 {
		t.Errorf("Expected 1 open forecast, got %d", len(openForecasts))
	}

	// Test listing resolved forecasts
	resolvedForecasts, err := service.ForecastList(ctx, "resolved", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(resolvedForecasts) != 1 {
		t.Errorf("Expected 1 resolved forecast, got %d", len(resolvedForecasts))
	}

	// Test invalid list type
	_, err = service.ForecastList(ctx, "invalid", "")
	if err == nil || err.Error() != "invalid resolved status" {
		t.Errorf("Expected 'invalid resolved status' error, got %v", err)
	}

	// Test error cases
	getForecastRepoMock(mockRepo).listOpenErr = errors.New("database error")
	_, err = service.ForecastList(ctx, "open", "")
	if err == nil {
		t.Error("Expected error from ListOpenForecasts, got nil")
	}

	getForecastRepoMock(mockRepo).listOpenErr = nil
	getForecastRepoMock(mockRepo).listResolvedErr = errors.New("database error")
	_, err = service.ForecastList(ctx, "resolved", "")
	if err == nil {
		t.Error("Expected error from ListResolvedForecasts, got nil")
	}
}
