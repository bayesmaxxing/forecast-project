package services

import (
	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"database/sql"
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

func (f *ForecastPointService) GetForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	cacheKey := fmt.Sprintf("point:list:%d", id)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}
	points, err := f.repo.GetForecastPointsByForecastID(ctx, id)

	if err != nil {
		return nil, err
	}
	f.cache.Set(cacheKey, points)

	return points, nil
}

func (f *ForecastPointService) GetForecastPointsByForecastIDAndUser(ctx context.Context, id int64, user_id int64) ([]*models.ForecastPoint, error) {
	cacheKey := fmt.Sprintf("point:list:user:%d:%d", user_id, id)

	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}
	points, err := f.repo.GetForecastPointsByForecastIDAndUser(ctx, id, user_id)
	if err != nil {
		return nil, err
	}
	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetOrderedForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	return f.repo.GetOrderedForecastPointsByForecastID(ctx, id)
}

func (f *ForecastPointService) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	// Check if the forecast exists
	_, err := f.f_repo.GetForecastByID(ctx, fp.ForecastID)
	if err == sql.ErrNoRows {
		return errors.New("forecast not found")
	}
	if err != nil {
		return err
	}

	f.cache.Delete(fmt.Sprintf("point:list:%d", fp.ForecastID))
	f.cache.Delete("point:all:latest")
	f.cache.Delete("point:all")

	return f.repo.CreateForecastPoint(ctx, fp)
}

func (f *ForecastPointService) GetAllForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	cacheKey := "point:all"
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetAllForecastPoints(ctx)
	if err != nil {
		return nil, err
	}

	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetLatestForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	cacheKey := "point:all:latest"
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetLatestForecastPoints(ctx)
	if err != nil {
		return nil, err
	}

	f.cache.Set(cacheKey, points)
	return points, nil
}

func (f *ForecastPointService) GetLatestForecastPointsByUser(ctx context.Context, user_id int64) ([]*models.ForecastPoint, error) {
	cacheKey := fmt.Sprintf("point:all:latest:%d", user_id)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]*models.ForecastPoint), nil
	}

	points, err := f.repo.GetLatestForecastPointsByUser(ctx, user_id)
	if err != nil {
		return nil, err
	}
	f.cache.Set(cacheKey, points)
	return points, nil
}

type GraphPoint struct {
	PointForecast float64   `json:"point_forecast"`
	CreatedAt     time.Time `json:"created"`
	UserID        int64     `json:"user_id"`
}

func (f *ForecastPointService) GetOrderedForecastPoints(ctx context.Context, forecastID int64, userID int64) ([]GraphPoint, error) {

	cacheKey := fmt.Sprintf("point:ordered:%d:%d", userID, forecastID)
	if cachedPoints, found := f.cache.Get(cacheKey); found {
		return cachedPoints.([]GraphPoint), nil
	}

	var points []*models.ForecastPoint

	var err error
	if userID != 0 {
		points, err = f.repo.GetForecastPointsByForecastIDAndUser(ctx, forecastID, userID)
	} else {
		// If no user_id is provided, return all points
		points, err = f.repo.GetOrderedForecastPointsByForecastID(ctx, forecastID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	graphPoints := make([]GraphPoint, len(points))
	for i, p := range points {
		graphPoints[i] = GraphPoint{
			PointForecast: p.PointForecast,
			CreatedAt:     p.CreatedAt,
			UserID:        p.UserID,
		}
	}

	f.cache.Set(cacheKey, graphPoints)
	return graphPoints, nil
}

func (f *ForecastPointService) GetTodaysForecastPoints(ctx context.Context, user_id int64) ([]*models.ForecastPoint, error) {

	points, err := f.repo.GetTodaysForecastPoints(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return points, nil
}
