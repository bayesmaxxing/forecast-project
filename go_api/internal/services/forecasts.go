package services

import (
	"context"
	"errors"
	"go_api/internal/models"
	"go_api/internal/repository"
	"go_api/internal/utils"
)

type ForecastService struct {
	repo      *repository.ForecastRepository
	pointRepo *repository.ForecastPointRepository
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

	probabilities := make([]float64, len(points))
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

func (s *ForecastService) GetAggregatedScores(ctx context.Context, category string) (*utils.AggScores, error) {
	var forecasts []*models.Forecast
	var err error

	if category == "" {
		forecasts, err = s.repo.ListResolvedWithScores(ctx)
	} else {
		forecasts, err = s.repo.ListResolvedWithScoresAndCategory(ctx, category)
	}

	if err != nil {
		return &utils.AggScores{}, err
	}

	scores := make([]utils.ForecastScores, 0, len(forecasts))
	for i, forecast := range forecasts {
		scores[i] = utils.ForecastScores{
			BrierScore: *forecast.BrierScore,
			Log2Score:  *forecast.Log2Score,
			LogNScore:  *forecast.LogNScore,
		}

	}

	aggScores, err := utils.CalculateAggregateScores(scores, category)
	if err != nil {
		return &utils.AggScores{}, err
	}

	return &aggScores, nil
}
