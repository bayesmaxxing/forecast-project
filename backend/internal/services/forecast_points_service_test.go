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

// MockForecastPointRepository implements the repository.ForecastPointRepository interface
type MockForecastPointRepository struct {
	points               []*models.ForecastPoint
	nextID               int64
	getAllErr            error
	getByForecastIDErr   error
	getByForecastUserErr error
	createErr            error
	getLatestErr         error
	getLatestByUserErr   error
}

func NewMockForecastPointRepository() repository.ForecastPointRepository {
	return &MockForecastPointRepository{
		points: make([]*models.ForecastPoint, 0),
		nextID: 1,
	}
}

func (m *MockForecastPointRepository) GetAllForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.points, nil
}

func (m *MockForecastPointRepository) GetForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	if m.getByForecastIDErr != nil {
		return nil, m.getByForecastIDErr
	}
	var result []*models.ForecastPoint
	for _, p := range m.points {
		if p.ForecastID == id {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockForecastPointRepository) GetForecastPointsByForecastIDAndUser(ctx context.Context, id int64, userID int64) ([]*models.ForecastPoint, error) {
	if m.getByForecastUserErr != nil {
		return nil, m.getByForecastUserErr
	}
	var result []*models.ForecastPoint
	for _, p := range m.points {
		if p.ForecastID == id && p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockForecastPointRepository) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	if m.createErr != nil {
		return m.createErr
	}
	fp.ID = m.nextID
	m.nextID++
	fp.CreatedAt = time.Now()
	m.points = append(m.points, fp)
	return nil
}

func (m *MockForecastPointRepository) GetLatestForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	if m.getLatestErr != nil {
		return nil, m.getLatestErr
	}
	// Simple implementation that doesn't actually filter for latest
	return m.points, nil
}

func (m *MockForecastPointRepository) GetLatestForecastPointsByUser(ctx context.Context, userID int64) ([]*models.ForecastPoint, error) {
	if m.getLatestByUserErr != nil {
		return nil, m.getLatestByUserErr
	}
	var result []*models.ForecastPoint
	for _, p := range m.points {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

// Helper function to access the underlying mock
func getMockForecastPointRepo(repo repository.ForecastPointRepository) *MockForecastPointRepository {
	return repo.(*MockForecastPointRepository)
}

// Test for GetAllForecastPoints
func TestGetAllForecastPoints(t *testing.T) {
	// Setup
	mockRepo := NewMockForecastPointRepository()
	service := services.NewForecastPointService(mockRepo)
	ctx := context.Background()

	// Create test forecast points
	point1 := &models.ForecastPoint{
		ForecastID:    1,
		PointForecast: 0.7,
		UserID:        1,
	}
	point2 := &models.ForecastPoint{
		ForecastID:    2,
		PointForecast: 0.8,
		UserID:        2,
	}
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point1)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point2)

	// Test successful retrieval
	points, err := service.GetAllForecastPoints(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(points) != 2 {
		t.Errorf("Expected 2 points, got %d", len(points))
	}

	// Test error case
	getMockForecastPointRepo(mockRepo).getAllErr = errors.New("database error")
	_, err = service.GetAllForecastPoints(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetForecastPointsByForecastID
func TestGetForecastPointsByForecastID(t *testing.T) {
	// Setup
	mockRepo := NewMockForecastPointRepository()
	service := services.NewForecastPointService(mockRepo)
	ctx := context.Background()

	// Create test forecast points
	forecastID := int64(1)
	point1 := &models.ForecastPoint{
		ForecastID:    forecastID,
		PointForecast: 0.7,
		UserID:        1,
	}
	point2 := &models.ForecastPoint{
		ForecastID:    forecastID,
		PointForecast: 0.8,
		UserID:        2,
	}
	point3 := &models.ForecastPoint{
		ForecastID:    2,
		PointForecast: 0.6,
		UserID:        1,
	}
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point1)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point2)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point3)

	// Test successful retrieval
	points, err := service.GetForecastPointsByForecastID(ctx, forecastID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(points) != 2 {
		t.Errorf("Expected 2 points for forecast ID %d, got %d", forecastID, len(points))
	}

	// Test error case
	getMockForecastPointRepo(mockRepo).getByForecastIDErr = errors.New("database error")
	_, err = service.GetForecastPointsByForecastID(ctx, forecastID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetForecastPointsByForecastIDAndUser
func TestGetForecastPointsByForecastIDAndUser(t *testing.T) {
	// Setup
	mockRepo := NewMockForecastPointRepository()
	service := services.NewForecastPointService(mockRepo)
	ctx := context.Background()

	// Create test forecast points
	forecastID := int64(1)
	userID := int64(1)
	point1 := &models.ForecastPoint{
		ForecastID:    forecastID,
		PointForecast: 0.7,
		UserID:        userID,
	}
	point2 := &models.ForecastPoint{
		ForecastID:    forecastID,
		PointForecast: 0.8,
		UserID:        2,
	}
	point3 := &models.ForecastPoint{
		ForecastID:    2,
		PointForecast: 0.6,
		UserID:        userID,
	}
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point1)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point2)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point3)

	// Test successful retrieval
	points, err := service.GetForecastPointsByForecastIDAndUser(ctx, forecastID, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(points) != 1 {
		t.Errorf("Expected 1 point for forecast ID %d and user ID %d, got %d", forecastID, userID, len(points))
	}

	// Test error case
	getMockForecastPointRepo(mockRepo).getByForecastUserErr = errors.New("database error")
	_, err = service.GetForecastPointsByForecastIDAndUser(ctx, forecastID, userID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for CreateForecastPoint
func TestCreateForecastPoint(t *testing.T) {
	// Setup
	mockRepo := NewMockForecastPointRepository()
	service := services.NewForecastPointService(mockRepo)
	ctx := context.Background()

	// Test successful creation
	testPoint := &models.ForecastPoint{
		ForecastID:    1,
		PointForecast: 0.7,
		UserID:        1,
		Reason:        "Test reason",
	}
	err := service.CreateForecastPoint(ctx, testPoint)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if testPoint.ID != 1 {
		t.Errorf("Expected point ID 1, got %d", testPoint.ID)
	}

	// Verify point was stored
	points := getMockForecastPointRepo(mockRepo).points
	if len(points) != 1 {
		t.Errorf("Expected 1 point stored, got %d", len(points))
	}
	if points[0].PointForecast != testPoint.PointForecast {
		t.Errorf("Expected point forecast %f, got %f", testPoint.PointForecast, points[0].PointForecast)
	}

	// Test error case
	getMockForecastPointRepo(mockRepo).createErr = errors.New("database error")
	err = service.CreateForecastPoint(ctx, testPoint)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetLatestForecastPoints
func TestGetLatestForecastPoints(t *testing.T) {
	// Setup
	mockRepo := NewMockForecastPointRepository()
	service := services.NewForecastPointService(mockRepo)
	ctx := context.Background()

	// Create test forecast points
	point1 := &models.ForecastPoint{
		ForecastID:    1,
		PointForecast: 0.7,
		UserID:        1,
	}
	point2 := &models.ForecastPoint{
		ForecastID:    2,
		PointForecast: 0.8,
		UserID:        2,
	}
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point1)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point2)

	// Test successful retrieval
	points, err := service.GetLatestForecastPoints(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(points) != 2 {
		t.Errorf("Expected 2 points, got %d", len(points))
	}

	// Test error case
	getMockForecastPointRepo(mockRepo).getLatestErr = errors.New("database error")
	_, err = service.GetLatestForecastPoints(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetLatestForecastPointsByUser
func TestGetLatestForecastPointsByUser(t *testing.T) {
	// Setup
	mockRepo := NewMockForecastPointRepository()
	service := services.NewForecastPointService(mockRepo)
	ctx := context.Background()

	// Create test forecast points
	userID := int64(1)
	point1 := &models.ForecastPoint{
		ForecastID:    1,
		PointForecast: 0.7,
		UserID:        userID,
	}
	point2 := &models.ForecastPoint{
		ForecastID:    2,
		PointForecast: 0.8,
		UserID:        userID,
	}
	point3 := &models.ForecastPoint{
		ForecastID:    3,
		PointForecast: 0.6,
		UserID:        2,
	}
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point1)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point2)
	getMockForecastPointRepo(mockRepo).CreateForecastPoint(ctx, point3)

	// Test successful retrieval
	points, err := service.GetLatestForecastPointsByUser(ctx, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(points) != 2 {
		t.Errorf("Expected 2 points for user ID %d, got %d", userID, len(points))
	}

	// Test error case
	getMockForecastPointRepo(mockRepo).getLatestByUserErr = errors.New("database error")
	_, err = service.GetLatestForecastPointsByUser(ctx, userID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
