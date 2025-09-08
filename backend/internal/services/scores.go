package services

import (
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"strconv"
)

type ScoreService struct {
	repo  repository.ScoreRepository
	cache *cache.Cache
}

func NewScoreService(repo repository.ScoreRepository, cache *cache.Cache) *ScoreService {
	return &ScoreService{repo: repo, cache: cache}
}

// multiple-score methods
func (s *ScoreService) GetScores(ctx context.Context, user_id int64, forecast_id int64) ([]models.Scores, error) {
	switch {
	case user_id != 0 && forecast_id != 0:
		cacheKey := fmt.Sprintf("score:by_user_and_forecast:%d:%d", user_id, forecast_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			return cachedData.([]models.Scores), nil
		}
		scores, err := s.repo.GetScoreByForecastAndUser(ctx, forecast_id, user_id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return []models.Scores{*scores}, nil

	case user_id != 0 && forecast_id == 0:
		cacheKey := fmt.Sprintf("score:by_user:%d", user_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			return cachedData.([]models.Scores), nil
		}
		scores, err := s.repo.GetScoresByUserID(ctx, user_id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	case forecast_id != 0 && user_id == 0:
		cacheKey := fmt.Sprintf("score:by_forecast:%d", forecast_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			return cachedData.([]models.Scores), nil
		}
		scores, err := s.repo.GetScoreByForecastID(ctx, forecast_id)
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	case user_id == 0 && forecast_id == 0:
		cacheKey := "score:all"
		if cachedData, found := s.cache.Get(cacheKey); found {
			return cachedData.([]models.Scores), nil
		}
		scores, err := s.repo.GetAllScores(ctx)
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	}
	return nil, errors.New("no scores found")
}

func (s *ScoreService) GetScoreByForecastID(ctx context.Context, forecast_id int64) ([]models.Scores, error) {
	return s.repo.GetScoreByForecastID(ctx, forecast_id)
}

func (s *ScoreService) GetScoreByForecastAndUser(ctx context.Context, forecast_id int64, user_id int64) (*models.Scores, error) {
	return s.repo.GetScoreByForecastAndUser(ctx, forecast_id, user_id)
}

func (s *ScoreService) GetAverageScoreByForecastID(ctx context.Context, forecast_id int64) (*models.ScoreMetrics, error) {
	cacheKey := "score:detail:average:" + strconv.FormatInt(forecast_id, 10)

	// Try to get from cache first
	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.ScoreMetrics), nil
	}
	score, err := s.repo.GetAverageScoreByForecastID(ctx, forecast_id)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, score)
	return score, nil
}

func (s *ScoreService) GetScoresByUserID(ctx context.Context, user_id int64) ([]models.Scores, error) {
	return s.repo.GetScoresByUserID(ctx, user_id)
}

// manipulate score model
func (s *ScoreService) CreateScore(ctx context.Context, score *models.Scores) error {
	s.cache.DeleteByPrefix("score:")

	return s.repo.CreateScore(ctx, score)
}

func (s *ScoreService) UpdateScore(ctx context.Context, score *models.Scores) error {
	return s.repo.UpdateScore(ctx, score)
}

func (s *ScoreService) DeleteScore(ctx context.Context, score_id int64) error {
	s.cache.DeleteByPrefix("score:")

	return s.repo.DeleteScore(ctx, score_id)
}

func (s *ScoreService) GetAllScores(ctx context.Context) ([]models.Scores, error) {
	cacheKey := "score:all"

	// Try to get from cache first
	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.([]models.Scores), nil
	}
	scores, err := s.repo.GetAllScores(ctx)
	if err != nil {
		return nil, err
	}

	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAverageScores(ctx context.Context) ([]models.Scores, error) {
	cacheKey := "score:all:average"

	// Try to get from cache first
	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.([]models.Scores), nil
	}
	scores, err := s.repo.GetAverageScores(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

// Aggregate Scores
func (s *ScoreService) GetAggregateScoresByUserID(ctx context.Context, user_id int64) (*models.UserScores, error) {

	cacheKey := fmt.Sprintf("score:aggregate:%d", user_id)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.UserScores), nil
	}

	scores, err := s.repo.GetUserOverallScores(ctx, user_id)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUserIDAndCategory(ctx context.Context, user_id int64, category string) (*models.UserCategoryScores, error) {
	cacheKey := fmt.Sprintf("score:aggregate:%d:%s", user_id, category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.UserCategoryScores), nil
	}
	scores, err := s.repo.GetUserCategoryScores(ctx, user_id, category)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUsers(ctx context.Context) ([]models.UserScores, error) {
	cacheKey := "score:aggregate:users"

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.([]models.UserScores), nil
	}
	scores, err := s.repo.GetOverallScoresByUsers(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByCategory(ctx context.Context, category string) (*models.CategoryScores, error) {
	cacheKey := fmt.Sprintf("score:aggregate:%s", category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.CategoryScores), nil
	}
	scores, err := s.repo.GetCategoryScores(ctx, category)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUsersAndCategory(ctx context.Context, category string) ([]models.UserCategoryScores, error) {
	cacheKey := fmt.Sprintf("score:aggregate:users:%s", category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.([]models.UserCategoryScores), nil
	}
	scores, err := s.repo.GetCategoryScoresByUsers(ctx, category)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

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
