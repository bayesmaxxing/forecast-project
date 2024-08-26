package services

import (
	"context"
	"errors"
	"go_api/internal/models"
	"go_api/internal/repository"
)

type ForecastService struct {
	repo      *repository.ForecastRepository
	pointRepo *repository.ForecastPointRepository
}

func NewForecastRepository(repo *repository.ForecastRepository, pointRepo *repository.ForecastPointRepository) *ForecastService {
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

func (s *ForecastService) DeleteForecast(ctx context.Context, id int64) error {
	return s.repo.DeleteForecast(ctx, id)
}

func (s *ForecastService) ResolveForecast(ctx context.Context,
	f *models.Forecast, id int64, resolution string, comment string) error {
	forecast, err := s.repo.GetForecastByID(ctx, id)
	if err != nil {
		return err
	}

	points, err := s.pointRepo.GetForecastPointsByForecastID(ctx, id)
	if err != nil {
		return err
	}

	var probabilities []float64

	probabilities = make([]float64, len(points))
	for i, point := range points {
		probabilities[i] = point.PointForecast
	}

	if err := forecast.Resolve(resolution, comment, probabilities); err != nil {
		return err
	}

	return s.repo.UpdateForecast(ctx, forecast)
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