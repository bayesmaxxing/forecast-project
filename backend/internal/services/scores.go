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
	"time"
)

type ScoreService struct {
	repo  repository.ScoreRepository
	cache *cache.Cache
}

func NewScoreService(repo repository.ScoreRepository, cache *cache.Cache) *ScoreService {
	return &ScoreService{repo: repo, cache: cache}
}

// getCacheableDateRangeKey returns a cache key suffix if the date range matches a predefined
// cacheable range (all_time, last_12_months, last_6_months). Returns empty string and false
// if the range is custom and should not be cached.
func getCacheableDateRangeKey(startDate *time.Time, endDate *time.Time) (string, bool) {
	// All time (no dates) - cacheable
	if startDate == nil && endDate == nil {
		return "all_time", true
	}

	// If end date is specified and not "now", it's a custom range
	if endDate != nil {
		now := time.Now()
		// Allow some tolerance (1 day) for "now"
		if endDate.Before(now.AddDate(0, 0, -1)) {
			return "", false
		}
	}

	if startDate == nil {
		return "", false
	}

	now := time.Now()

	// Check if it's approximately 12 months ago (allow 1 day tolerance)
	twelveMonthsAgo := now.AddDate(-1, 0, 0)
	if startDate.After(twelveMonthsAgo.AddDate(0, 0, -1)) && startDate.Before(twelveMonthsAgo.AddDate(0, 0, 2)) {
		return "last_12_months", true
	}

	// Check if it's approximately 6 months ago
	sixMonthsAgo := now.AddDate(0, -6, 0)
	if startDate.After(sixMonthsAgo.AddDate(0, 0, -1)) && startDate.Before(sixMonthsAgo.AddDate(0, 0, 2)) {
		return "last_6_months", true
	}

	// Check if it's approximately 3 months ago
	threeMonthsAgo := now.AddDate(0, -3, 0)
	if startDate.After(threeMonthsAgo.AddDate(0, 0, -1)) && startDate.Before(threeMonthsAgo.AddDate(0, 0, 2)) {
		return "last_3_months", true
	}

	// Custom range - not cacheable
	return "", false
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
func (s *ScoreService) GetAggregateScores(ctx context.Context, user_id *int64, forecast_id *int64, category *string, startDate *time.Time, endDate *time.Time) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores", slog.Any("user_id", user_id), slog.Any("forecast_id", forecast_id), slog.Any("category", category), slog.Any("start_date", startDate), slog.Any("end_date", endDate))
	switch {
	case user_id != nil && category != nil:
		log.Info("getting aggregate scores by user and category", slog.Any("user_id", user_id), slog.Any("category", category))
		return s.GetAggregateScoresByUserIDAndCategory(ctx, models.ScoreFilters{UserID: user_id, Category: category, StartDate: startDate, EndDate: endDate})
	case user_id != nil:
		log.Info("getting aggregate scores by user", slog.Any("user_id", user_id))
		return s.GetAggregateScoresByUserID(ctx, models.ScoreFilters{UserID: user_id, StartDate: startDate, EndDate: endDate})
	case category != nil:
		log.Info("getting aggregate scores by category", slog.Any("category", category))
		return s.GetAggregateScoresByCategory(ctx, models.ScoreFilters{Category: category, StartDate: startDate, EndDate: endDate})
	case forecast_id != nil:
		log.Info("getting aggregate scores by forecast", slog.Any("forecast_id", forecast_id))
		return s.GetAggregateScoresByForecastID(ctx, models.ScoreFilters{ForecastID: forecast_id, StartDate: startDate, EndDate: endDate})
	default:
		log.Info("getting overall scores")
		return s.GetOverallScores(ctx, models.ScoreFilters{StartDate: startDate, EndDate: endDate})
	}
}

func (s *ScoreService) GetAggregateScoresByUserID(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)
	log.Info("getting aggregate scores by user", slog.Any("user_id", filters.UserID))

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:%d:%s", *filters.UserID, dateRangeKey)
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

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache", slog.Any("user_id", filters.UserID))
	return s.repo.GetAggregateScores(ctx, filters)
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

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:%d:%s:%s", userID, category, dateRangeKey)
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

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetAggregateScores(ctx, filters)
}

func (s *ScoreService) GetAggregateScoresByForecastID(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)
	log.Info("getting aggregate scores by forecast", slog.Any("forecast_id", filters.ForecastID))

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:forecast:%d:%s", *filters.ForecastID, dateRangeKey)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by forecast"))
			return cachedData.(*models.OverallScores), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "aggregate scores by forecast"))
		score, err := s.repo.GetAggregateScores(ctx, filters)
		if err != nil {
			log.Error("failed to get aggregate scores by forecast", slog.Any("forecast_id", filters.ForecastID), slog.String("error", err.Error()))
			return nil, err
		}
		s.cache.Set(cacheKey, score)
		return score, nil
	}

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetAggregateScores(ctx, filters)
}

func (s *ScoreService) GetAggregateScoresByCategory(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)
	log.Info("getting aggregate scores by category", slog.Any("category", filters.Category))

	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:%s:%s", category, dateRangeKey)
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

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetAggregateScores(ctx, filters)
}

func (s *ScoreService) GetOverallScores(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	log := logger.FromContext(ctx)
	log.Info("getting overall scores")

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:overall:%s", dateRangeKey)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "overall scores"))
			return cachedData.(*models.OverallScores), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "overall scores"))
		scores, err := s.repo.GetAggregateScores(ctx, filters)
		if err != nil {
			log.Error("failed to get overall scores", slog.String("error", err.Error()))
			return nil, err
		}
		s.cache.Set(cacheKey, scores)
		return scores, nil
	}

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetAggregateScores(ctx, filters)
}

// router for group by user aggregate scores
func (s *ScoreService) GetAggregateScoresGroupedByUsers(ctx context.Context, category *string, startDate *time.Time, endDate *time.Time) ([]models.UserScores, error) {
	log := logger.FromContext(ctx)

	log.Info("getting aggregate scores grouped by users", slog.Any("category", category))
	groupByUserID := true
	if category != nil {
		log.Info("getting aggregate scores grouped by users and category", slog.Any("category", category))
		return s.GetAggregateScoresByUsersAndCategory(ctx, models.ScoreFilters{Category: category, GroupByUserID: &groupByUserID, StartDate: startDate, EndDate: endDate})
	} else {
		log.Info("getting aggregate scores grouped by users")
		return s.GetAggregateScoresByUsers(ctx, models.ScoreFilters{GroupByUserID: &groupByUserID, StartDate: startDate, EndDate: endDate})
	}
}

func (s *ScoreService) GetAggregateScoresByUsers(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {
	log := logger.FromContext(ctx)
	log.Info("getting aggregate scores grouped by users", slog.Any("filters", filters))

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:users:%s", dateRangeKey)
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

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetAggregateScoresByUsers(ctx, filters)
}

func (s *ScoreService) GetAggregateScoresByUsersAndCategory(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {
	log := logger.FromContext(ctx)

	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}

	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := fmt.Sprintf("score:aggregate:users:%s:%s", category, dateRangeKey)
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

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetAggregateScoresByUsers(ctx, filters)
}
