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

type ForecastPointService struct {
	repo   repository.ForecastPointRepository
	f_repo repository.ForecastRepository
	cache  *cache.Cache
}

func NewForecastPointService(fp_repo repository.ForecastPointRepository, f_repo repository.ForecastRepository, cache *cache.Cache) *ForecastPointService {
	return &ForecastPointService{repo: fp_repo, f_repo: f_repo, cache: cache}
}

// routes handler requests to the associated service method based on filters
// this is to ensure that caching is handled correctly
func (f *ForecastPointService) GetForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {

	log := logger.FromContext(ctx)
	log.Info("routing to forecast point service", slog.Any("filters", filters))
	switch {
	// user and forecast id are provided
	case filters.UserID != nil && filters.ForecastID != nil:
		log.Info("user and forecast id are provided, fetching by forecast id and user id")
		return f.GetForecastPointsByForecastIDAndUser(ctx, filters)

	case filters.UserID != nil && filters.Date != nil:
		log.Info("user and date are provided, fetching by date")
		return f.GetForecastPointsByDate(ctx, filters)

	case filters.UserID != nil && filters.ForecastID == nil:
		log.Info("user is provided, fetching latest forecast points by user")
		return f.GetLatestForecastPointsByUser(ctx, filters)

	case filters.UserID == nil && filters.ForecastID != nil:
		log.Info("forecast id is provided, fetching by forecast id")
		return f.GetForecastPointsByForecastID(ctx, filters)

	case filters.UserID == nil:
		log.Info("no user or forecast id is provided, fetching latest forecast points")
		return f.GetLatestForecastPoints(ctx, filters)

	case filters.UserID == nil && filters.ForecastID == nil:
		log.Info("no user or forecast id is provided, fetching all forecast points")
		return f.GetAllForecastPoints(ctx, filters)
	}
	log.Info("no route found, fetching forecast points", slog.Any("filters", filters))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		log.Error("failed to fetch forecast points", slog.String("error", err.Error()))
		return nil, err
	}
	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByForecastID(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	log := logger.FromContext(ctx)
	log.Info("fetching forecast points by forecast id", slog.Any("filters", filters))
	cacheKey := fmt.Sprintf("point:list:%d", filters.ForecastID)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "forecast points by forecast id"))
		return cachedPoints.([]*models.ForecastPoint), nil
	}
	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "forecast points by forecast id"))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		log.Error("failed to fetch forecast points from database", slog.String("error", err.Error()))
		return nil, err
	}
	f.cache.Set(cacheKey, points)

	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByForecastIDAndUser(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	log := logger.FromContext(ctx)

	cacheKey := fmt.Sprintf("point:list:user:%d:%d", filters.UserID, filters.ForecastID)

	if cachedPoints, found := f.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "forecast points by forecast id and user id"))
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "forecast points by forecast id and user id"))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		log.Error("failed to fetch forecast points from database", slog.String("error", err.Error()))
		return nil, err
	}
	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	log := logger.FromContext(ctx)
	// Check if the forecast exists
	log.Info("checking if forecast exists", slog.Int64("forecast_id", fp.ForecastID))
	forecast, err := f.f_repo.GetForecastByID(ctx, fp.ForecastID)
	if err != nil {
		log.Error("failed to get forecast", slog.String("error", err.Error()))
		return err
	}
	if forecast == nil {
		log.Error("forecast not found")
		return errors.New("forecast not found")
	}

	// Check if forecast is already resolved
	if forecast.ResolvedAt != nil {
		log.Error("forecast has already been resolved")
		return errors.New("forecast has already been resolved")
	}

	// Check if forecast closing date has passed
	if forecast.ClosingDate != nil && forecast.ClosingDate.Before(time.Now()) {
		log.Error("forecast has already closed")
		return errors.New("forecast has already closed")
	}

	log.Info("deleting cache keys",
		slog.String("cache_key", fmt.Sprintf("point:list:%d", fp.ForecastID)),
		slog.String("cache_key", "point:all:latest"),
		slog.String("cache_key", "point:all"))
	f.cache.Delete(fmt.Sprintf("point:list:%d", fp.ForecastID))
	f.cache.Delete("point:all:latest")
	f.cache.Delete("point:all")

	log.Info("creating forecast point", slog.Any("forecast_point", fp))
	return f.repo.CreateForecastPoint(ctx, fp)
}

func (f *ForecastPointService) GetAllForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	log := logger.FromContext(ctx)

	cacheKey := "point:all"
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "all forecast points"))
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "all forecast points"))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		log.Error("failed to fetch forecast points from database", slog.String("error", err.Error()))
		return nil, err
	}

	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetLatestForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	log := logger.FromContext(ctx)

	cacheKey := "point:all:latest"
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "latest forecast points"))
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "latest forecast points"))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		log.Error("failed to fetch forecast points from database", slog.String("error", err.Error()))
		return nil, err
	}

	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetLatestForecastPointsByUser(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	log := logger.FromContext(ctx)

	cacheKey := fmt.Sprintf("point:all:latest:%d", filters.UserID)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "latest forecast points by user"))
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "latest forecast points by user"))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		log.Error("failed to fetch forecast points from database", slog.String("error", err.Error()))
		return nil, err
	}
	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByDate(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	log := logger.FromContext(ctx)

	log.Info("fetching latest forecast points by user", slog.Any("filters", filters))
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}
	return points, nil
}
