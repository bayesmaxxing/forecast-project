package services

import (
	"backend/internal/cache"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	log := logger.FromContext(ctx)
	switch {
	case user_id != 0 && forecast_id != 0:
		log.Info("getting scores by user and forecast", slog.Int64("user_id", user_id), slog.Int64("forecast_id", forecast_id))
		cacheKey := fmt.Sprintf("score:by_user_and_forecast:%d:%d", user_id, forecast_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "scores by user and forecast"))
			return cachedData.([]models.Scores), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "scores by user and forecast"))
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{UserID: &user_id, ForecastID: &forecast_id})
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil

	case user_id != 0 && forecast_id == 0:
		log.Info("getting scores by user", slog.Int64("user_id", user_id))
		cacheKey := fmt.Sprintf("score:by_user:%d", user_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "scores by user"))
			return cachedData.([]models.Scores), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "scores by user"))
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{UserID: &user_id})
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	case forecast_id != 0 && user_id == 0:
		log.Info("getting scores by forecast", slog.Int64("forecast_id", forecast_id))
		cacheKey := fmt.Sprintf("score:by_forecast:%d", forecast_id)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "scores by forecast"))
			return cachedData.([]models.Scores), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "scores by forecast"))
		scores, err := s.repo.GetScores(ctx, models.ScoreFilters{ForecastID: &forecast_id})
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	case user_id == 0 && forecast_id == 0:
		log.Info("getting all scores")
		cacheKey := "score:all"
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "all scores"))
			return cachedData.([]models.Scores), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "all scores"))
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
	log := logger.FromContext(ctx)

	log.Info("getting average scores")
	cacheKey := "score:all:average"

	// Try to get from cache first
	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "average scores"))
		return cachedData.([]models.Scores), nil
	}
	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "average scores"))
	scores, err := s.repo.GetAverageScores(ctx)
	if err != nil {
		log.Error("failed to get average scores", slog.String("error", err.Error()))
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

// Aggregate Scores router
func (s *ScoreService) GetAggregateScores(ctx context.Context, user_id *int64, forecast_id *int64, category *string) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores", slog.Any("user_id", user_id), slog.Any("forecast_id", forecast_id), slog.Any("category", category))
	switch {
	case user_id != nil && category != nil:
		log.Info("getting aggregate scores by user and category", slog.Any("user_id", user_id), slog.Any("category", category))
		return s.GetAggregateScoresByUserIDAndCategory(ctx, models.ScoreFilters{UserID: user_id, Category: category})
	case user_id != nil:
		log.Info("getting aggregate scores by user", slog.Any("user_id", user_id))
		return s.GetAggregateScoresByUserID(ctx, models.ScoreFilters{UserID: user_id})
	case category != nil:
		log.Info("getting aggregate scores by category", slog.Any("category", category))
		return s.GetAggregateScoresByCategory(ctx, models.ScoreFilters{Category: category})
	case forecast_id != nil:
		log.Info("getting aggregate scores by forecast", slog.Any("forecast_id", forecast_id))
		return s.GetAggregateScoresByForecastID(ctx, models.ScoreFilters{ForecastID: forecast_id})
	default:
		log.Info("getting overall scores")
		return s.GetOverallScores(ctx)
	}
}

func (s *ScoreService) GetAggregateScoresByUserID(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)
	cacheKey := fmt.Sprintf("score:aggregate:%d", filters.UserID)
	log.Info("getting aggregate scores by user", slog.Any("user_id", filters.UserID))
	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by user"))
		return cachedData.(*models.OverallScores), nil
	}

	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by user"))
	scores, err := s.repo.GetAggregateScores(ctx, filters)
	if err != nil {
		log.Error("failed to get aggregate scores by user", slog.Any("user_id", filters.UserID), slog.String("error", err.Error()))
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUserIDAndCategory(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores by user and category", slog.Any("user_id", filters.UserID), slog.Any("category", filters.Category))
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
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by user and category"))
		return cachedData.(*models.OverallScores), nil
	}

	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by user and category"))
	scores, err := s.repo.GetAggregateScores(ctx, filters)
	if err != nil {
		log.Error("failed to get aggregate scores by user and category", slog.Any("user_id", filters.UserID), slog.Any("category", filters.Category), slog.String("error", err.Error()))
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByForecastID(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {

	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores by forecast", slog.Any("forecast_id", filters.ForecastID))
	cacheKey := fmt.Sprintf("score:aggregate:%d", filters.ForecastID)

	// Try to get from cache first
	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by forecast"))
		return cachedData.(*models.OverallScores), nil
	}

	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by forecast"))
	score, err := s.repo.GetAggregateScores(ctx, models.ScoreFilters{ForecastID: filters.ForecastID})
	if err != nil {
		log.Error("failed to get aggregate scores by forecast", slog.Any("forecast_id", filters.ForecastID), slog.String("error", err.Error()))
		return nil, err
	}
	s.cache.Set(cacheKey, score)
	return score, nil
}

func (s *ScoreService) GetAggregateScoresByCategory(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {

	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores by category", slog.Any("category", filters.Category))
	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}
	cacheKey := fmt.Sprintf("score:aggregate:%s", category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by category"))
		return cachedData.(*models.OverallScores), nil
	}

	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by category"))
	scores, err := s.repo.GetAggregateScores(ctx, filters)
	if err != nil {
		log.Error("failed to get aggregate scores by category", slog.Any("category", filters.Category), slog.String("error", err.Error()))
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting overall scores")
	cacheKey := "score:aggregate:overall"
	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "overall scores"))
		return cachedData.(*models.OverallScores), nil
	}
	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "overall scores"))
	scores, err := s.repo.GetAggregateScores(ctx, models.ScoreFilters{})
	if err != nil {
		log.Error("failed to get overall scores", slog.String("error", err.Error()))
		return nil, err
	}

	s.cache.Set(cacheKey, scores)
	return scores, nil
}

// router for group by user aggregate scores
func (s *ScoreService) GetAggregateScoresGroupedByUsers(ctx context.Context, category *string) ([]models.UserScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores grouped by users", slog.Any("category", category))
	groupByUserID := true
	if category != nil {
		log.Info("getting aggregate scores grouped by users and category", slog.Any("category", category))
		return s.GetAggregateScoresByUsersAndCategory(ctx, models.ScoreFilters{Category: category, GroupByUserID: &groupByUserID})
	} else {
		log.Info("getting aggregate scores grouped by users")
		return s.GetAggregateScoresByUsers(ctx, models.ScoreFilters{GroupByUserID: &groupByUserID})
	}
}

func (s *ScoreService) GetAggregateScoresByUsers(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores grouped by users", slog.Any("filters", filters))
	cacheKey := "score:aggregate:users"

	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores grouped by users"))
		return cachedData.([]models.UserScores), nil
	}

	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores grouped by users"))
	scores, err := s.repo.GetAggregateScoresByUsers(ctx, filters)
	if err != nil {
		log.Error("failed to get aggregate scores grouped by users", slog.Any("filters", filters), slog.String("error", err.Error()))
		return nil, err
	}
	s.cache.Set(cacheKey, scores)
	return scores, nil
}

func (s *ScoreService) GetAggregateScoresByUsersAndCategory(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {

	log := logger.FromContext(ctx)

	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}
	cacheKey := fmt.Sprintf("score:aggregate:users:%s", category)

	if cachedData, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores grouped by users and category"))
		return cachedData.([]models.UserScores), nil
	}
	log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores grouped by users and category"))
	scores, err := s.repo.GetAggregateScoresByUsers(ctx, filters)
	if err != nil {
		log.Error("failed to get aggregate scores grouped by users and category", slog.Any("filters", filters), slog.String("error", err.Error()))
		return nil, err
	}

	s.cache.Set(cacheKey, scores)
	return scores, nil
}
