package services

import (
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
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

	switch {
	// user and forecast id are provided
	case filters.UserID != nil && filters.ForecastID != nil:
		return f.GetForecastPointsByForecastIDAndUser(ctx, filters)

	case filters.UserID != nil && filters.Date != nil:
		return f.GetForecastPointsByDate(ctx, filters)

	case filters.UserID != nil && filters.ForecastID == nil:
		return f.GetLatestForecastPointsByUser(ctx, filters)

	case filters.UserID == nil && filters.ForecastID != nil:
		return f.GetForecastPointsByForecastID(ctx, filters)

	case filters.UserID == nil:
		return f.GetLatestForecastPoints(ctx, filters)

	case filters.UserID == nil && filters.ForecastID == nil:
		return f.GetAllForecastPoints(ctx, filters)
	}
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByForecastID(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	cacheKey := fmt.Sprintf("point:list:%d", filters.ForecastID)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}
	points, err := f.repo.GetForecastPoints(ctx, filters)

	if err != nil {
		return nil, err
	}
	f.cache.Set(cacheKey, points)

	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByForecastIDAndUser(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	cacheKey := fmt.Sprintf("point:list:user:%d:%d", filters.UserID, filters.ForecastID)

	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}
	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	// Check if the forecast exists
	forecast, err := f.f_repo.GetForecastByID(ctx, fp.ForecastID)
	if err != nil {
		return err
	}
	if forecast == nil {
		return errors.New("forecast not found")
	}

	// Check if forecast is already resolved
	if forecast.ResolvedAt != nil {
		return errors.New("forecast has already been resolved")
	}

	// Check if forecast closing date has passed
	if forecast.ClosingDate != nil && forecast.ClosingDate.Before(time.Now()) {
		return errors.New("forecast has already closed")
	}

	f.cache.Delete(fmt.Sprintf("point:list:%d", fp.ForecastID))
	f.cache.Delete("point:all:latest")
	f.cache.Delete("point:all")

	return f.repo.CreateForecastPoint(ctx, fp)
}

func (f *ForecastPointService) GetAllForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	cacheKey := "point:all"
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}

	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetLatestForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	cacheKey := "point:all:latest"
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}

	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetLatestForecastPointsByUser(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	cacheKey := fmt.Sprintf("point:all:latest:%d", filters.UserID)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}
	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByDate(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	points, err := f.repo.GetForecastPoints(ctx, filters)
	if err != nil {
		return nil, err
	}
	return points, nil
}
