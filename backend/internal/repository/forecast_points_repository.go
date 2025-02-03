package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"time"

	_ "github.com/jackc/pgx/v5"
)

type ForecastPointRepository struct {
	db *database.DB
}

func NewForecastPointRepository(db *database.DB) *ForecastPointRepository {
	return &ForecastPointRepository{db: db}
}

func (r *ForecastPointRepository) GetAllForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				update_id
				, forecast_id
				, point_forecast
				, upper_ci
				, lower_ci
				, reason
				, created
				, user_id
				FROM forecast_points `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecast_points []*models.ForecastPoint
	for rows.Next() {
		var fp models.ForecastPoint

		if err := rows.Scan(&fp.ID,
			&fp.ForecastID,
			&fp.PointForecast,
			&fp.UpperCI,
			&fp.LowerCI,
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *ForecastPointRepository) GetForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				update_id
				, forecast_id
				, point_forecast
				, upper_ci
				, lower_ci
				, reason
				, created
				, user_id
				FROM forecast_points 
				WHERE forecast_id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecast_points []*models.ForecastPoint
	for rows.Next() {
		var fp models.ForecastPoint
		if err := rows.Scan(&fp.ID,
			&fp.ForecastID,
			&fp.PointForecast,
			&fp.UpperCI,
			&fp.LowerCI,
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *ForecastPointRepository) GetForecastPointsByForecastIDAndUser(ctx context.Context, id int64, user_id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				update_id
				, forecast_id
				, point_forecast
				, upper_ci
				, lower_ci
				, reason
				, created
				, user_id
				FROM forecast_points 
				WHERE forecast_id = $1
				AND user_id = $2`

	rows, err := r.db.QueryContext(ctx, query, id, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecast_points []*models.ForecastPoint
	for rows.Next() {
		var fp models.ForecastPoint
		if err := rows.Scan(&fp.ID,
			&fp.ForecastID,
			&fp.PointForecast,
			&fp.UpperCI,
			&fp.LowerCI,
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *ForecastPointRepository) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	fp.CreatedAt = time.Now()

	query := `INSERT INTO forecast_points (forecast_id
											, point_forecast
											, upper_ci
											, lower_ci
											, created
											, reason
											, user_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING update_id`

	err := r.db.QueryRowContext(ctx, query, fp.ForecastID, fp.PointForecast, fp.UpperCI,
		fp.LowerCI, fp.CreatedAt, fp.Reason, fp.UserID).Scan(&fp.ID)
	return err
}

func (r *ForecastPointRepository) GetLatestForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	query := `SELECT distinct on (forecast_id)
				update_id
				, forecast_id
				, point_forecast
				, upper_ci
				, lower_ci
				, created
				, reason
				, user_id
				FROM forecast_points
				ORDER BY forecast_id, created DESC;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecast_points []*models.ForecastPoint
	for rows.Next() {
		var fp models.ForecastPoint
		if err := rows.Scan(&fp.ID,
			&fp.ForecastID,
			&fp.PointForecast,
			&fp.UpperCI,
			&fp.LowerCI,
			&fp.CreatedAt,
			&fp.Reason,
			&fp.UserID); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *ForecastPointRepository) GetLatestForecastPointsByUser(ctx context.Context, user_id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT distinct on (forecast_id)
				update_id
				, forecast_id
				, point_forecast
				, upper_ci
				, lower_ci
				, created
				, reason
				, user_id
				FROM forecast_points
				WHERE user_id = $1
				ORDER BY forecast_id, created DESC;`

	rows, err := r.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecast_points []*models.ForecastPoint
	for rows.Next() {
		var fp models.ForecastPoint
		if err := rows.Scan(&fp.ID,
			&fp.ForecastID,
			&fp.PointForecast,
			&fp.UpperCI,
			&fp.LowerCI,
			&fp.CreatedAt,
			&fp.Reason,
			&fp.UserID); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}
