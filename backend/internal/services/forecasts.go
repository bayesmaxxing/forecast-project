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
	cacheKey := fmt.Sprintf("forecast:detail:%d", id)

	if cachedForecast, found := s.cache.Get(cacheKey); found {
		return cachedForecast.(*models.Forecast), nil
	}

	forecast, err := s.repo.GetForecastByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.cache.Set(cacheKey, forecast)

	return forecast, nil
}

func (s *ForecastService) CheckForecastOwnership(ctx context.Context, id int64, user_id int64) (bool, error) {
	return s.repo.CheckForecastOwnership(ctx, id, user_id)
}

func (s *ForecastService) CheckForecastStatus(ctx context.Context, id int64) (bool, error) {
	return s.repo.CheckForecastStatus(ctx, id)
}

func (s *ForecastService) CreateForecast(ctx context.Context, f *models.Forecast) error {
	s.cache.DeleteByPrefix("forecast:list:")

	return s.repo.CreateForecast(ctx, f)
}

func (s *ForecastService) DeleteForecast(ctx context.Context, id int64, user_id int64) error {
	s.cache.DeleteByPrefix("forecast:list:")

	return s.repo.DeleteForecast(ctx, id, user_id)
}

func (s *ForecastService) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	s.cache.DeleteByPrefix("forecast:list:")

	return s.repo.UpdateForecast(ctx, f)
}

func (s *ForecastService) ResolveForecast(ctx context.Context, user_id int64, id int64, resolution string, comment string) error {
	forecast, err := s.repo.GetForecastByID(ctx, id)
	if err != nil {
		return err
	}

	// Make sure the forecast exists and the user owns it
	ownership, err := s.CheckForecastOwnership(ctx, id, user_id)
	if err == sql.ErrNoRows {
		return errors.New("forecast does not exist")
	}

	if !ownership {
		return errors.New("user does not own this forecast")
	}

	// Make sure the forecast is not already resolved
	status, err := s.CheckForecastStatus(ctx, id)
	if err != nil {
		return err
	}

	if !status {
		return errors.New("forecast is already resolved")
	}

	points, err := s.pointRepo.GetForecastPointsByForecastID(ctx, id)
	if err != nil {
		return err
	}

	if len(points) == 0 {
		return errors.New("no forecast points found")
	}

	// Group points by user
	userPoints := make(map[int64][]float64)
	for _, point := range points {
		userPoints[point.UserID] = append(userPoints[point.UserID], point.PointForecast)
	}

	// Update the forecast with resolution status
	now := time.Now()
	forecast.ResolvedAt = &now
	forecast.Resolution = &resolution
	forecast.ResolutionComment = &comment

	// Update the forecast in the database
	if err := s.repo.UpdateForecast(ctx, forecast); err != nil {
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

		score, err := models.CalcForecastScore(probabilities, outcome, userID, forecast.ID)
		if err != nil {
			return err
		}

		if err := s.scoreRepo.CreateScore(ctx, &score); err != nil {
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
func (s *ForecastService) ForecastList(ctx context.Context, listType string, category string) ([]*models.Forecast, error) {
	cacheKey := fmt.Sprintf("forecast:list:%s:%s", listType, category)

	if cachedList, found := s.cache.Get(cacheKey); found {
		return cachedList.([]*models.Forecast), nil
	}

	switch listType {
	case "open":
		if category != "" {
			forecasts, err := s.repo.ListOpenForecastsWithCategory(ctx, category)
			if err != nil {
				return nil, err
			}
			s.cache.Set(cacheKey, forecasts)
			return forecasts, nil
		}
		forecasts, err := s.repo.ListOpenForecasts(ctx)
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, forecasts)
		return forecasts, nil
	case "resolved":
		if category != "" {
			forecasts, err := s.repo.ListResolvedForecastsWithCategory(ctx, category)
			if err != nil {
				return nil, err
			}
			s.cache.Set(cacheKey, forecasts)
			return forecasts, nil
		}
		forecasts, err := s.repo.ListResolvedForecasts(ctx)
		if err != nil {
			return nil, err
		}
		s.cache.Set(cacheKey, forecasts)
		return forecasts, nil
	default:
		return nil, errors.New("invalid resolved status")
	}
}

func (s *ForecastService) GetStaleAndNewForecasts(ctx context.Context, userID int64) ([]*models.Forecast, error) {
	forecasts, err := s.repo.GetStaleAndNewForecasts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return forecasts, nil
}
