package services

import (
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
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
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{UserID: &user_id, ForecastID: &forecast_id})
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil

	case user_id != 0 && forecast_id == 0:
		cacheKey := fmt.Sprintf("score:by_user:%d", user_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			return cachedData.([]models.Scores), nil
		}
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{UserID: &user_id})
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
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{ForecastID: &forecast_id})
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
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{})
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	}
	return nil, errors.New("no scores found")
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

// Aggregate Scores router
func (s *ScoreService) GetAggregateScores(ctx context.Context, user_id *int64, forecast_id *int64, category *string) (*models.OverallScores, error) {
	switch {
	case user_id != nil && category != nil:
		return s.GetAggregateScoresByUserIDAndCategory(ctx, models.ScoreFilters{UserID: user_id, Category: category})
	case user_id != nil:
		return s.GetAggregateScoresByUserID(ctx, models.ScoreFilters{UserID: user_id})
	case category != nil:
		return s.GetAggregateScoresByCategory(ctx, models.ScoreFilters{Category: category})
	case forecast_id != nil:
		return s.GetAggregateScoresByForecastID(ctx, models.ScoreFilters{ForecastID: forecast_id})
	default:
		return s.GetOverallScores(ctx)
	}
}

func (s *ScoreService) GetAggregateScoresByUserID(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {

	cacheKey := fmt.Sprintf("score:aggregate:%d", filters.UserID)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.OverallScores), nil
	}

	scores, err := s.repo.GetAggregateScores(ctx, filters)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUserIDAndCategory(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	userID := int64(0)
	if filters.UserID != nil {
		userID = *filters.UserID
	}
	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}
	cacheKey := fmt.Sprintf("score:aggregate:%d:%s", userID, category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.OverallScores), nil
	}
	scores, err := s.repo.GetAggregateScores(ctx, filters)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByForecastID(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	cacheKey := fmt.Sprintf("score:aggregate:%d", filters.ForecastID)

	// Try to get from cache first
	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.OverallScores), nil
	}
	score, err := s.repo.GetAggregateScores(ctx, models.ScoreFilters{ForecastID: filters.ForecastID})
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, score)
	return score, nil
}

func (s *ScoreService) GetAggregateScoresByCategory(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}
	cacheKey := fmt.Sprintf("score:aggregate:%s", category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.OverallScores), nil
	}
	scores, err := s.repo.GetAggregateScores(ctx, filters)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	cacheKey := "score:aggregate:overall"
	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.(*models.OverallScores), nil
	}
	scores, err := s.repo.GetAggregateScores(ctx, models.ScoreFilters{})
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

// router for group by user aggregate scores
func (s *ScoreService) GetAggregateScoresGroupedByUsers(ctx context.Context, category *string) ([]models.UserScores, error) {
	groupByUserID := true
	if category != nil {
		return s.GetAggregateScoresByUsersAndCategory(ctx, models.ScoreFilters{Category: category, GroupByUserID: &groupByUserID})
	} else {
		return s.GetAggregateScoresByUsers(ctx, models.ScoreFilters{GroupByUserID: &groupByUserID})
	}
}

func (s *ScoreService) GetAggregateScoresByUsers(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {
	cacheKey := "score:aggregate:users"

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.([]models.UserScores), nil
	}
	scores, err := s.repo.GetAggregateScoresByUsers(ctx, filters)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUsersAndCategory(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {
	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}
	cacheKey := fmt.Sprintf("score:aggregate:users:%s", category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		return cachedData.([]models.UserScores), nil
	}
	scores, err := s.repo.GetAggregateScoresByUsers(ctx, filters)
	if err != nil {
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}
