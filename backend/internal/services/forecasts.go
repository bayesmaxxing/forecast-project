package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
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

	probabilities := make([]float64, 0, len(points))
	for _, point := range points {
		probabilities = append(probabilities, point.PointForecast)
	}
	log.Printf("Probabilities supplied: %f", probabilities[0])
	if err := forecast.Resolve(resolution, comment, probabilities); err != nil {
		return err
	}
	log.Printf("Forecast resolved")
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

	if len(forecasts) == 0 {
		return &utils.AggScores{
			AggBrierScore: 0,
			AggLog2Score:  0,
			AggLogNScore:  0,
			Category:      category,
		}, nil
	}

	scores := make([]utils.ForecastScores, 0, len(forecasts))
	for _, forecast := range forecasts {
		if forecast.BrierScore == nil {
			continue
		}
		scores = append(scores, utils.ForecastScores{
			BrierScore: *forecast.BrierScore,
			Log2Score:  *forecast.Log2Score,
			LogNScore:  *forecast.LogNScore,
		})
	}
	if len(scores) == 0 {
		return &utils.AggScores{
			AggBrierScore: 0,
			AggLog2Score:  0,
			AggLogNScore:  0,
			Category:      category,
		}, nil
	}

	aggScores, err := utils.CalculateAggregateScores(scores, category)
	if err != nil {
		return nil, err
	}

	return &aggScores, nil
}
