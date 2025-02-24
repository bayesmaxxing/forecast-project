package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type ForecastPointService struct {
	repo repository.ForecastPointRepository
}

func NewForecastPointService(repo repository.ForecastPointRepository) *ForecastPointService {
	return &ForecastPointService{repo: repo}
}

func (f *ForecastPointService) GetForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	return f.repo.GetForecastPointsByForecastID(ctx, id)
}

func (f *ForecastPointService) GetForecastPointsByForecastIDAndUser(ctx context.Context, id int64, user_id int64) ([]*models.ForecastPoint, error) {
	return f.repo.GetForecastPointsByForecastIDAndUser(ctx, id, user_id)
}

func (f *ForecastPointService) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	return f.repo.CreateForecastPoint(ctx, fp)
}

func (f *ForecastPointService) GetAllForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	return f.repo.GetAllForecastPoints(ctx)
}

func (f *ForecastPointService) GetLatestForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	return f.repo.GetLatestForecastPoints(ctx)
}

func (f *ForecastPointService) GetLatestForecastPointsByUser(ctx context.Context, user_id int64) ([]*models.ForecastPoint, error) {
	return f.repo.GetLatestForecastPointsByUser(ctx, user_id)
}
