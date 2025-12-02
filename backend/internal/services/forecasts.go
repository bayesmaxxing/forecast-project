package services

import (
	"backend/internal/cache"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type ForecastService struct {
	repo      repository.ForecastRepository
	pointRepo repository.ForecastPointRepository
	scoreRepo repository.ScoreRepository
	cache     *cache.Cache
}

func NewForecastService(repo repository.ForecastRepository, pointRepo repository.ForecastPointRepository, scoreRepo repository.ScoreRepository, cache *cache.Cache) *ForecastService {
	return &ForecastService{
		repo:      repo,
		pointRepo: pointRepo,
		scoreRepo: scoreRepo,
		cache:     cache,
	}
}

// individual forecast operations
func (s *ForecastService) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
	log := logger.FromContext(ctx)

	cacheKey := fmt.Sprintf("forecast:detail:%d", id)

	if cachedForecast, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "forecast by id"))
		return cachedForecast.(*models.Forecast), nil
	}

	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "forecast by id"))
	forecast, err := s.repo.GetForecastByID(ctx, id)
	if err != nil {
		log.Error("failed to fetch forecast from database", slog.String("error", err.Error()))
		return nil, err
	}

	s.cache.Set(cacheKey, forecast)

	return forecast, nil
}

func (s *ForecastService) CheckForecastOwnership(ctx context.Context, id int64, user_id int64) (bool, error) {
	log := logger.FromContext(ctx)

	log.Info("checking forecast ownership", slog.Int64("id", id), slog.Int64("user_id", user_id))
	return s.repo.CheckForecastOwnership(ctx, id, user_id)
}

func (s *ForecastService) CheckForecastStatus(ctx context.Context, id int64) (bool, error) {
	log := logger.FromContext(ctx)

	log.Info("checking forecast status", slog.Int64("id", id))
	return s.repo.CheckForecastStatus(ctx, id)
}

func (s *ForecastService) CreateForecast(ctx context.Context, f *models.Forecast) error {
	log := logger.FromContext(ctx)

	log.Info("creating forecast", slog.Any("forecast", f))
	s.cache.DeleteByPrefix("forecast:list:")

	return s.repo.CreateForecast(ctx, f)
}

func (s *ForecastService) DeleteForecast(ctx context.Context, id int64, user_id int64) error {
	log := logger.FromContext(ctx)

	log.Info("deleting forecast", slog.Int64("id", id), slog.Int64("user_id", user_id))
	s.cache.DeleteByPrefix("forecast:list:")

	return s.repo.DeleteForecast(ctx, id, user_id)
}

func (s *ForecastService) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	log := logger.FromContext(ctx)

	log.Info("updating forecast", slog.Any("forecast", f))
	s.cache.DeleteByPrefix("forecast:list:")

	return s.repo.UpdateForecast(ctx, f)
}

func (s *ForecastService) ResolveForecast(ctx context.Context, user_id int64, id int64, resolution string, comment string) error {
	log := logger.FromContext(ctx)

	log.Info("resolving forecast", slog.Int64("id", id), slog.Int64("user_id", user_id), slog.String("resolution", resolution), slog.String("comment", comment))
	forecast, err := s.repo.GetForecastByID(ctx, id)
	if err != nil {
		log.Error("failed to get forecast from database", slog.String("error", err.Error()))
		return err
	}

	// Make sure the forecast exists and the user owns it
	log.Info("checking forecast ownership", slog.Int64("id", id), slog.Int64("user_id", user_id))
	ownership, err := s.CheckForecastOwnership(ctx, id, user_id)
	if err == sql.ErrNoRows {
		log.Error("forecast does not exist", slog.Int64("id", id), slog.Int64("user_id", user_id))
		return errors.New("forecast does not exist")
	}

	if !ownership {
		log.Error("user does not own this forecast", slog.Int64("id", id), slog.Int64("user_id", user_id))
		return errors.New("user does not own this forecast")
	}

	// Make sure the forecast is not already resolved
	status, err := s.CheckForecastStatus(ctx, id)
	if err != nil {
		log.Error("failed to check forecast status", slog.Int64("id", id), slog.Int64("user_id", user_id), slog.String("error", err.Error()))
		return err
	}

	if !status {
		log.Error("forecast is already resolved", slog.Int64("id", id), slog.Int64("user_id", user_id))
		return errors.New("forecast is already resolved")
	}

	points, err := s.pointRepo.GetForecastPoints(ctx, models.PointFilters{ForecastID: &id})
	if err != nil {
		log.Error("failed to get forecast points", slog.Int64("id", id), slog.Int64("user_id", user_id), slog.String("error", err.Error()))
		return err
	}

	if len(points) == 0 {
		log.Error("no forecast points found", slog.Int64("id", id), slog.Int64("user_id", user_id))
		return errors.New("no forecast points found")
	}

	// Group points and created at by user
	userPoints := make(map[int64][]models.TimePoint)
	for _, point := range points {
		userPoints[point.UserID] = append(userPoints[point.UserID], models.TimePoint{
			PointForecast: point.PointForecast,
			CreatedAt:     point.CreatedAt,
		})
	}

	// Update the forecast with resolution status
	now := time.Now()
	forecast.ResolvedAt = &now
	forecast.Resolution = &resolution
	forecast.ResolutionComment = &comment

	// Update the forecast in the database
	if err := s.repo.UpdateForecast(ctx, forecast); err != nil {
		log.Error("failed to update forecast in database", slog.Int64("id", id), slog.Int64("user_id", user_id), slog.String("error", err.Error()))
		return err
	}

	if resolution == "-" {
		return nil
	}

	outcome := resolution == "1"
	for userID, probabilities := range userPoints {
		if len(probabilities) == 0 {
			continue
		}
		log.Info("calculating forecast score")
		score, err := models.CalcForecastScore(probabilities, outcome, userID, forecast.ID, forecast.CreatedAt, forecast.ClosingDate, forecast.ResolvedAt)
		if err != nil {
			log.Error("failed to calculate forecast score", slog.Int64("id", id), slog.Int64("user_id", user_id), slog.String("error", err.Error()))
			return err
		}

		if err := s.scoreRepo.CreateScore(ctx, &score); err != nil {
			log.Error("failed to create score", slog.Int64("id", id), slog.Int64("user_id", user_id), slog.String("error", err.Error()))
			return err
		}
	}

	// invalidate affected cache keys
	deleteKey := fmt.Sprintf("forecast:detail:%d", forecast.ID)
	s.cache.Delete(deleteKey)
	s.cache.DeleteByPrefix("forecast:list:")
	s.cache.DeleteByPrefix("scores:")

	return nil
}

// aggregate forecast operations
func (s *ForecastService) GetForecasts(ctx context.Context, filters models.ForecastFilters) ([]*models.Forecast, error) {
	log := logger.FromContext(ctx)

	log.Info("getting forecasts", slog.Any("filters", filters))
	status := ""
	if filters.Status != nil {
		status = *filters.Status
	}

	category := ""
	if filters.Category != nil {
		category = *filters.Category
	}

	cacheKey := fmt.Sprintf("forecast:list:%s:%s", status, category)

	if cachedList, found := s.cache.Get(cacheKey); found {
		log.Info("cache hit",
			slog.String("cache_key", cacheKey),
			slog.String("cache_type", "forecasts"))
		return cachedList.([]*models.Forecast), nil
	}

	if filters.ForecastID != nil {
		return nil, errors.New("forecast ID is not supported for this operation")
	}

	log.Info("cache miss",
		slog.String("cache_key", cacheKey),
		slog.String("cache_type", "forecasts"))
	forecasts, err := s.repo.GetForecasts(ctx, filters)
	if err != nil {
		return nil, err
	}

	s.cache.Set(cacheKey, forecasts)
	return forecasts, nil
}

func (s *ForecastService) GetStaleAndNewForecasts(ctx context.Context, userID int64) ([]*models.Forecast, error) {
	log := logger.FromContext(ctx)

	log.Info("getting stale and new forecasts", slog.Int64("user_id", userID))
	forecasts, err := s.repo.GetStaleAndNewForecasts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return forecasts, nil
}
