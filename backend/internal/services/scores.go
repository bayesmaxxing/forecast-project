package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type ScoreService struct {
	repo repository.ScoreRepository
}

func NewScoreService(repo repository.ScoreRepository) *ScoreService {
	return &ScoreService{repo: repo}
}

// single-score methods
func (s *ScoreService) GetScoreByForecastID(ctx context.Context, forecast_id int64) ([]models.Scores, error) {
	return s.repo.GetScoreByForecastID(ctx, forecast_id)
}

func (s *ScoreService) GetScoreByForecastAndUser(ctx context.Context, forecast_id int64, user_id int64) (*models.Scores, error) {
	return s.repo.GetScoreByForecastAndUser(ctx, forecast_id, user_id)
}

func (s *ScoreService) GetAverageScoreByForecastID(ctx context.Context, forecast_id int64) (*models.ScoreMetrics, error) {
	return s.repo.GetAverageScoreByForecastID(ctx, forecast_id)
}

func (s *ScoreService) GetScoresByUserID(ctx context.Context, user_id int64) ([]models.Scores, error) {
	return s.repo.GetScoresByUserID(ctx, user_id)
}

// manipulate score model
func (s *ScoreService) CreateScore(ctx context.Context, score *models.Scores) error {
	return s.repo.CreateScore(ctx, score)
}

func (s *ScoreService) UpdateScore(ctx context.Context, score *models.Scores) error {
	return s.repo.UpdateScore(ctx, score)
}

func (s *ScoreService) DeleteScore(ctx context.Context, score_id int64) error {
	return s.repo.DeleteScore(ctx, score_id)
}

func (s *ScoreService) GetAllScores(ctx context.Context) ([]models.Scores, error) {
	return s.repo.GetAllScores(ctx)
}

// Aggregate Scores
func (s *ScoreService) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	return s.repo.GetOverallScores(ctx)
}

func (s *ScoreService) GetCategoryScores(ctx context.Context, category string) (*models.CategoryScores, error) {
	return s.repo.GetCategoryScores(ctx, category)
}

func (s *ScoreService) GetCategoryScoresByUsers(ctx context.Context, category string) ([]models.UserCategoryScores, error) {
	return s.repo.GetCategoryScoresByUsers(ctx, category)
}

func (s *ScoreService) GetOverallScoresByUsers(ctx context.Context) ([]models.UserScores, error) {
	return s.repo.GetOverallScoresByUsers(ctx)
}

// user-specific aggregate scores
func (s *ScoreService) GetUserOverallScores(ctx context.Context, user_id int64) (*models.UserScores, error) {
	return s.repo.GetUserOverallScores(ctx, user_id)
}

func (s *ScoreService) GetUserCategoryScores(ctx context.Context, user_id int64, category string) (*models.UserCategoryScores, error) {
	return s.repo.GetUserCategoryScores(ctx, user_id, category)
}
