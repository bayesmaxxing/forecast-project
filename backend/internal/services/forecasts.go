package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"time"
)

type ForecastService struct {
	repo      *repository.ForecastRepository
	pointRepo *repository.ForecastPointRepository
	scoreRepo *repository.ScoreRepository
}

func NewForecastService(repo *repository.ForecastRepository, pointRepo *repository.ForecastPointRepository) *ForecastService {
	return &ForecastService{
		repo:      repo,
		pointRepo: pointRepo,
	}
}

func (s *ForecastService) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
	return s.repo.GetForecastByID(ctx, id)
}

func (s *ForecastService) CreateForecast(ctx context.Context, f *models.Forecast) error {
	return s.repo.CreateForecast(ctx, f)
}

func (s *ForecastService) DeleteForecast(ctx context.Context, id int64, user_id int64) error {
	return s.repo.DeleteForecast(ctx, id, user_id)
}

func (s *ForecastService) UpdateForecast(ctx context.Context, f *models.Forecast, user_id int64) error {
	return s.repo.UpdateForecast(ctx, f, user_id)
}

func (s *ForecastService) ResolveForecast(ctx context.Context,
	id int64, resolution string, comment string) error {
	forecast, err := s.repo.GetForecastByID(ctx, id)
	if err != nil {
		return err
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
	if err := s.repo.UpdateForecast(ctx, forecast, forecast.UserID); err != nil {
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

	return nil
}

func (s *ForecastService) ForecastList(ctx context.Context, listType string, category string) ([]*models.Forecast, error) {
	switch listType {
	case "open":
		if category != "" {
			return s.repo.ListOpenForecastsWithCategory(ctx, category)
		}
		return s.repo.ListOpenForecasts(ctx)
	case "resolved":
		if category != "" {
			return s.repo.ListResolvedForecastsWithCategory(ctx, category)
		}
		return s.repo.ListResolvedForecasts(ctx)
	default:
		return nil, errors.New("invalid resolved status")
	}
}
