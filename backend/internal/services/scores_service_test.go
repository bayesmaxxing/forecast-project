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

// MockScoreRepository implements the repository.ScoreRepository interface
type MockScoreRepository struct {
	scores                  []models.Scores
	nextID                  int64
	getByForecastIDErr      error
	getByForecastUserErr    error
	getByUserIDErr          error
	createScoreErr          error
	updateScoreErr          error
	deleteScoreErr          error
	getAllScoresErr         error
	getOverallScoresErr     error
	getCategoryScoresErr    error
	getCatScoresByUserErr   error
	getOverallByUserErr     error
	getUserCatScoresErr     error
	getUserOverallScoresErr error
}

func NewMockScoreRepository() repository.ScoreRepository {
	return &MockScoreRepository{
		scores: make([]models.Scores, 0),
		nextID: 1,
	}
}

func (m *MockScoreRepository) GetScoreByForecastID(ctx context.Context, forecastID int64) ([]models.Scores, error) {
	if m.getByForecastIDErr != nil {
		return nil, m.getByForecastIDErr
	}
	var result []models.Scores
	for _, s := range m.scores {
		if s.ForecastID == forecastID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockScoreRepository) GetScoreByForecastAndUser(ctx context.Context, forecastID int64, userID int64) (*models.Scores, error) {
	if m.getByForecastUserErr != nil {
		return nil, m.getByForecastUserErr
	}
	for _, s := range m.scores {
		if s.ForecastID == forecastID && s.UserID == userID {
			return &s, nil
		}
	}
	return nil, errors.New("score not found")
}

func (m *MockScoreRepository) GetScoresByUserID(ctx context.Context, userID int64) ([]models.Scores, error) {
	if m.getByUserIDErr != nil {
		return nil, m.getByUserIDErr
	}
	var result []models.Scores
	for _, s := range m.scores {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockScoreRepository) CreateScore(ctx context.Context, score *models.Scores) error {
	if m.createScoreErr != nil {
		return m.createScoreErr
	}
	score.ID = m.nextID
	m.nextID++
	score.CreatedAt = time.Now()
	m.scores = append(m.scores, *score)
	return nil
}

func (m *MockScoreRepository) UpdateScore(ctx context.Context, score *models.Scores) error {
	if m.updateScoreErr != nil {
		return m.updateScoreErr
	}
	for i, s := range m.scores {
		if s.ID == score.ID {
			m.scores[i] = *score
			return nil
		}
	}
	return errors.New("score not found")
}

func (m *MockScoreRepository) DeleteScore(ctx context.Context, scoreID int64) error {
	if m.deleteScoreErr != nil {
		return m.deleteScoreErr
	}
	for i, s := range m.scores {
		if s.ID == scoreID {
			m.scores = append(m.scores[:i], m.scores[i+1:]...)
			return nil
		}
	}
	return errors.New("score not found")
}

func (m *MockScoreRepository) GetAllScores(ctx context.Context) ([]models.Scores, error) {
	if m.getAllScoresErr != nil {
		return nil, m.getAllScoresErr
	}
	return m.scores, nil
}

func (m *MockScoreRepository) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	if m.getOverallScoresErr != nil {
		return nil, m.getOverallScoresErr
	}
	// Simple mock implementation
	return &models.OverallScores{
		ScoreMetrics: models.ScoreMetrics{
			BrierScore: 0.25,
			Log2Score:  0.5,
			LogNScore:  0.75,
		},
		TotalUsers:     2,
		TotalForecasts: 3,
	}, nil
}

func (m *MockScoreRepository) GetCategoryScores(ctx context.Context, category string) (*models.CategoryScores, error) {
	if m.getCategoryScoresErr != nil {
		return nil, m.getCategoryScoresErr
	}
	// Simple mock implementation
	return &models.CategoryScores{
		ScoreMetrics: models.ScoreMetrics{
			BrierScore: 0.3,
			Log2Score:  0.6,
			LogNScore:  0.9,
		},
		TotalUsers:     2,
		TotalForecasts: 2,
		Category:       category,
	}, nil
}

func (m *MockScoreRepository) GetCategoryScoresByUsers(ctx context.Context, category string) ([]models.UserCategoryScores, error) {
	if m.getCatScoresByUserErr != nil {
		return nil, m.getCatScoresByUserErr
	}
	// Simple mock implementation
	return []models.UserCategoryScores{
		{
			ScoreMetrics: models.ScoreMetrics{
				BrierScore: 0.3,
				Log2Score:  0.6,
				LogNScore:  0.9,
			},
			UserID:         1,
			Category:       category,
			TotalForecasts: 2,
		},
		{
			ScoreMetrics: models.ScoreMetrics{
				BrierScore: 0.4,
				Log2Score:  0.7,
				LogNScore:  1.0,
			},
			UserID:         2,
			Category:       category,
			TotalForecasts: 1,
		},
	}, nil
}

func (m *MockScoreRepository) GetOverallScoresByUsers(ctx context.Context) ([]models.UserScores, error) {
	if m.getOverallByUserErr != nil {
		return nil, m.getOverallByUserErr
	}
	// Simple mock implementation
	return []models.UserScores{
		{
			ScoreMetrics: models.ScoreMetrics{
				BrierScore: 0.25,
				Log2Score:  0.5,
				LogNScore:  0.75,
			},
			UserID:         1,
			TotalForecasts: 2,
		},
		{
			ScoreMetrics: models.ScoreMetrics{
				BrierScore: 0.35,
				Log2Score:  0.65,
				LogNScore:  0.95,
			},
			UserID:         2,
			TotalForecasts: 1,
		},
	}, nil
}

func (m *MockScoreRepository) GetUserCategoryScores(ctx context.Context, userID int64, category string) (*models.UserCategoryScores, error) {
	if m.getUserCatScoresErr != nil {
		return nil, m.getUserCatScoresErr
	}
	// Simple mock implementation
	return &models.UserCategoryScores{
		ScoreMetrics: models.ScoreMetrics{
			BrierScore: 0.3,
			Log2Score:  0.6,
			LogNScore:  0.9,
		},
		UserID:         userID,
		Category:       category,
		TotalForecasts: 2,
	}, nil
}

func (m *MockScoreRepository) GetUserOverallScores(ctx context.Context, userID int64) (*models.UserScores, error) {
	if m.getUserOverallScoresErr != nil {
		return nil, m.getUserOverallScoresErr
	}
	// Simple mock implementation
	return &models.UserScores{
		ScoreMetrics: models.ScoreMetrics{
			BrierScore: 0.25,
			Log2Score:  0.5,
			LogNScore:  0.75,
		},
		UserID:         userID,
		TotalForecasts: 2,
	}, nil
}

// Helper function to access the underlying mock
func getMockScoreRepository(repo repository.ScoreRepository) *MockScoreRepository {
	return repo.(*MockScoreRepository)
}

// Test for GetScoreByForecastID
func TestGetScoreByForecastID(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Create test scores
	forecastID := int64(1)
	score1 := &models.Scores{
		ForecastID: forecastID,
		UserID:     1,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	score2 := &models.Scores{
		ForecastID: forecastID,
		UserID:     2,
		BrierScore: 0.3,
		Log2Score:  0.5,
		LogNScore:  0.7,
	}
	score3 := &models.Scores{
		ForecastID: 2,
		UserID:     1,
		BrierScore: 0.1,
		Log2Score:  0.3,
		LogNScore:  0.5,
	}
	getMockScoreRepository(mockRepo).CreateScore(ctx, score1)
	getMockScoreRepository(mockRepo).CreateScore(ctx, score2)
	getMockScoreRepository(mockRepo).CreateScore(ctx, score3)

	// Test successful retrieval
	scores, err := service.GetScoreByForecastID(ctx, forecastID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(scores) != 2 {
		t.Errorf("Expected 2 scores for forecast ID %d, got %d", forecastID, len(scores))
	}

	// Test error case
	getMockScoreRepository(mockRepo).getByForecastIDErr = errors.New("database error")
	_, err = service.GetScoreByForecastID(ctx, forecastID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetScoreByForecastAndUser
func TestGetScoreByForecastAndUser(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Create test scores
	forecastID := int64(1)
	userID := int64(1)
	score1 := &models.Scores{
		ForecastID: forecastID,
		UserID:     userID,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	score2 := &models.Scores{
		ForecastID: forecastID,
		UserID:     2,
		BrierScore: 0.3,
		Log2Score:  0.5,
		LogNScore:  0.7,
	}
	getMockScoreRepository(mockRepo).CreateScore(ctx, score1)
	getMockScoreRepository(mockRepo).CreateScore(ctx, score2)

	// Test successful retrieval
	score, err := service.GetScoreByForecastAndUser(ctx, forecastID, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if score.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, score.UserID)
	}
	if score.ForecastID != forecastID {
		t.Errorf("Expected forecast ID %d, got %d", forecastID, score.ForecastID)
	}

	// Test score not found
	_, err = service.GetScoreByForecastAndUser(ctx, 999, userID)
	if err == nil {
		t.Error("Expected error for nonexistent score, got nil")
	}

	// Test error case
	getMockScoreRepository(mockRepo).getByForecastUserErr = errors.New("database error")
	_, err = service.GetScoreByForecastAndUser(ctx, forecastID, userID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetScoresByUserID
func TestGetScoresByUserID(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Create test scores
	userID := int64(1)
	score1 := &models.Scores{
		ForecastID: 1,
		UserID:     userID,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	score2 := &models.Scores{
		ForecastID: 2,
		UserID:     userID,
		BrierScore: 0.3,
		Log2Score:  0.5,
		LogNScore:  0.7,
	}
	score3 := &models.Scores{
		ForecastID: 3,
		UserID:     2,
		BrierScore: 0.1,
		Log2Score:  0.3,
		LogNScore:  0.5,
	}
	getMockScoreRepository(mockRepo).CreateScore(ctx, score1)
	getMockScoreRepository(mockRepo).CreateScore(ctx, score2)
	getMockScoreRepository(mockRepo).CreateScore(ctx, score3)

	// Test successful retrieval
	scores, err := service.GetScoresByUserID(ctx, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(scores) != 2 {
		t.Errorf("Expected 2 scores for user ID %d, got %d", userID, len(scores))
	}

	// Test error case
	getMockScoreRepository(mockRepo).getByUserIDErr = errors.New("database error")
	_, err = service.GetScoresByUserID(ctx, userID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for CreateScore
func TestCreateScore(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Test successful creation
	testScore := &models.Scores{
		ForecastID: 1,
		UserID:     1,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	err := service.CreateScore(ctx, testScore)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if testScore.ID != 1 {
		t.Errorf("Expected score ID 1, got %d", testScore.ID)
	}

	// Verify score was stored
	scores := getMockScoreRepository(mockRepo).scores
	if len(scores) != 1 {
		t.Errorf("Expected 1 score stored, got %d", len(scores))
	}
	if scores[0].BrierScore != testScore.BrierScore {
		t.Errorf("Expected Brier score %f, got %f", testScore.BrierScore, scores[0].BrierScore)
	}

	// Test error case
	getMockScoreRepository(mockRepo).createScoreErr = errors.New("database error")
	err = service.CreateScore(ctx, testScore)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for UpdateScore
func TestUpdateScore(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Create a test score
	testScore := &models.Scores{
		ForecastID: 1,
		UserID:     1,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	getMockScoreRepository(mockRepo).CreateScore(ctx, testScore)

	// Update the score
	testScore.BrierScore = 0.3
	err := service.UpdateScore(ctx, testScore)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify score was updated
	scores := getMockScoreRepository(mockRepo).scores
	if scores[0].BrierScore != 0.3 {
		t.Errorf("Expected updated Brier score 0.3, got %f", scores[0].BrierScore)
	}

	// Test error case
	getMockScoreRepository(mockRepo).updateScoreErr = errors.New("database error")
	err = service.UpdateScore(ctx, testScore)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for DeleteScore
func TestDeleteScore(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Create a test score
	testScore := &models.Scores{
		ForecastID: 1,
		UserID:     1,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	getMockScoreRepository(mockRepo).CreateScore(ctx, testScore)

	// Delete the score
	err := service.DeleteScore(ctx, testScore.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify score was deleted
	scores := getMockScoreRepository(mockRepo).scores
	if len(scores) != 0 {
		t.Errorf("Expected 0 scores after deletion, got %d", len(scores))
	}

	// Test error case
	getMockScoreRepository(mockRepo).deleteScoreErr = errors.New("database error")
	err = service.DeleteScore(ctx, testScore.ID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetAllScores
func TestGetAllScores(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Create test scores
	score1 := &models.Scores{
		ForecastID: 1,
		UserID:     1,
		BrierScore: 0.2,
		Log2Score:  0.4,
		LogNScore:  0.6,
	}
	score2 := &models.Scores{
		ForecastID: 2,
		UserID:     2,
		BrierScore: 0.3,
		Log2Score:  0.5,
		LogNScore:  0.7,
	}
	getMockScoreRepository(mockRepo).CreateScore(ctx, score1)
	getMockScoreRepository(mockRepo).CreateScore(ctx, score2)

	// Test successful retrieval
	scores, err := service.GetAllScores(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(scores) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(scores))
	}

	// Test error case
	getMockScoreRepository(mockRepo).getAllScoresErr = errors.New("database error")
	_, err = service.GetAllScores(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for aggregate score methods
func TestAggregateScoreMethods(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()

	// Test GetOverallScores
	overallScores, err := service.GetOverallScores(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if overallScores.BrierScore != 0.25 {
		t.Errorf("Expected Brier score 0.25, got %f", overallScores.BrierScore)
	}

	// Test GetOverallScores error
	getMockScoreRepository(mockRepo).getOverallScoresErr = errors.New("database error")
	_, err = service.GetOverallScores(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Reset error
	getMockScoreRepository(mockRepo).getOverallScoresErr = nil

	// Test GetCategoryScores
	categoryScores, err := service.GetCategoryScores(ctx, "test-category")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if categoryScores.BrierScore != 0.3 {
		t.Errorf("Expected Brier score 0.3, got %f", categoryScores.BrierScore)
	}
	if categoryScores.Category != "test-category" {
		t.Errorf("Expected category 'test-category', got '%s'", categoryScores.Category)
	}

	// Test GetCategoryScores error
	getMockScoreRepository(mockRepo).getCategoryScoresErr = errors.New("database error")
	_, err = service.GetCategoryScores(ctx, "test-category")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Reset error
	getMockScoreRepository(mockRepo).getCategoryScoresErr = nil

	// Test GetCategoryScoresByUsers
	userCategoryScores, err := service.GetCategoryScoresByUsers(ctx, "test-category")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(userCategoryScores) != 2 {
		t.Errorf("Expected 2 user category scores, got %d", len(userCategoryScores))
	}
	if userCategoryScores[0].UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", userCategoryScores[0].UserID)
	}

	// Test GetCategoryScoresByUsers error
	getMockScoreRepository(mockRepo).getCatScoresByUserErr = errors.New("database error")
	_, err = service.GetCategoryScoresByUsers(ctx, "test-category")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Reset error
	getMockScoreRepository(mockRepo).getCatScoresByUserErr = nil

	// Test GetOverallScoresByUsers
	userScores, err := service.GetOverallScoresByUsers(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(userScores) != 2 {
		t.Errorf("Expected 2 user scores, got %d", len(userScores))
	}
	if userScores[0].UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", userScores[0].UserID)
	}

	// Test GetOverallScoresByUsers error
	getMockScoreRepository(mockRepo).getOverallByUserErr = errors.New("database error")
	_, err = service.GetOverallScoresByUsers(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for user-specific aggregate score methods
func TestUserAggregateScoreMethods(t *testing.T) {
	// Setup
	mockRepo := NewMockScoreRepository()
	service := services.NewScoreService(mockRepo)
	ctx := context.Background()
	userID := int64(1)

	// Test GetUserCategoryScores
	userCategoryScores, err := service.GetUserCategoryScores(ctx, userID, "test-category")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if userCategoryScores.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, userCategoryScores.UserID)
	}
	if userCategoryScores.Category != "test-category" {
		t.Errorf("Expected category 'test-category', got '%s'", userCategoryScores.Category)
	}

	// Test GetUserCategoryScores error
	getMockScoreRepository(mockRepo).getUserCatScoresErr = errors.New("database error")
	_, err = service.GetUserCategoryScores(ctx, userID, "test-category")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Reset error
	getMockScoreRepository(mockRepo).getUserCatScoresErr = nil

	// Test GetUserOverallScores
	userOverallScores, err := service.GetUserOverallScores(ctx, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if userOverallScores.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, userOverallScores.UserID)
	}

	// Test GetUserOverallScores error
	getMockScoreRepository(mockRepo).getUserOverallScoresErr = errors.New("database error")
	_, err = service.GetUserOverallScores(ctx, userID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
