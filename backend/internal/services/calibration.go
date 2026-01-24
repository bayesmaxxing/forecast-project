package services

import (
	"backend/internal/cache"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"log/slog"
)

type CalibrationService struct {
	repo  repository.CalibrationRepository
	cache *cache.Cache
}

func NewCalibrationService(repo repository.CalibrationRepository, cache *cache.Cache) *CalibrationService {
	return &CalibrationService{repo: repo, cache: cache}
}

func (s *CalibrationService) GetCalibrationData(ctx context.Context, filters models.CalibrationFilters) (*models.CalibrationData, error) {
	log := logger.FromContext(ctx)
	log.Info("getting calibration data", slog.Any("filters", filters))

	// Build cache key
	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := s.buildCacheKey("calibration", filters, dateRangeKey)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "calibration data"))
			return cachedData.(*models.CalibrationData), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "calibration data"))

		data, err := s.repo.GetCalibrationData(ctx, filters)
		if err != nil {
			log.Error("failed to get calibration data", slog.String("error", err.Error()))
			return nil, err
		}
		s.cache.Set(cacheKey, data)
		return data, nil
	}

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetCalibrationData(ctx, filters)
}

func (s *CalibrationService) GetCalibrationDataByUsers(ctx context.Context, filters models.CalibrationFilters) ([]models.UserCalibrationData, error) {
	log := logger.FromContext(ctx)
	log.Info("getting calibration data by users", slog.Any("filters", filters))

	// Build cache key
	dateRangeKey, cacheable := getCacheableDateRangeKey(filters.StartDate, filters.EndDate)
	if cacheable {
		cacheKey := s.buildCacheKey("calibration:users", filters, dateRangeKey)
		if cachedData, found := s.cache.Get(cacheKey); found {
			log.Info("cache hit", slog.String("cache_key", cacheKey), slog.String("cache_type", "calibration data by users"))
			return cachedData.([]models.UserCalibrationData), nil
		}
		log.Info("cache miss", slog.String("cache_key", cacheKey), slog.String("cache_type", "calibration data by users"))

		data, err := s.repo.GetCalibrationDataByUsers(ctx, filters)
		if err != nil {
			log.Error("failed to get calibration data by users", slog.String("error", err.Error()))
			return nil, err
		}
		s.cache.Set(cacheKey, data)
		return data, nil
	}

	// Custom date range - don't cache
	log.Info("custom date range - skipping cache")
	return s.repo.GetCalibrationDataByUsers(ctx, filters)
}

func (s *CalibrationService) buildCacheKey(prefix string, filters models.CalibrationFilters, dateRangeKey string) string {
	key := prefix

	if filters.UserID != nil {
		key = fmt.Sprintf("%s:user:%d", key, *filters.UserID)
	}
	if filters.Category != nil {
		key = fmt.Sprintf("%s:category:%s", key, *filters.Category)
	}
	key = fmt.Sprintf("%s:%s", key, dateRangeKey)

	return key
}
