package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"time"

	_ "github.com/jackc/pgx/v5"
)

type ForecastRepository struct {
	db *database.DB
}

func NewForecastRepository(db *database.DB) *ForecastRepository {
	return &ForecastRepository{db: db}
}

// Methods without user_id filtering
func (r *ForecastRepository) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
	query := `SELECT 
					id
					, question
					, category
					, created
					, user_id
					, resolution_criteria
					, resolution
					, resolved
					, comment
				FROM forecast_v2 
				WHERE id = $1`

	var f models.Forecast
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.Question,
		&f.Category,
		&f.CreatedAt,
		&f.ResolutionCriteria,
		&f.Resolution,
		&f.ResolvedAt,
		&f.ResolutionComment)

	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *ForecastRepository) CheckForecastOwnership(ctx context.Context, id int64, user_id int64) (bool, error) {
	query := `SELECT user_id FROM forecast_v2 WHERE id = $1`
	var forecastUserID int64
	err := r.db.QueryRowContext(ctx, query, id).Scan(&forecastUserID)
	if err != nil {
		return false, err
	}

	return forecastUserID == user_id, nil
}

func (r *ForecastRepository) CreateForecast(ctx context.Context, f *models.Forecast) error {
	f.CreatedAt = time.Now()

	query := `INSERT INTO forecast_v2 (
				question
				, category
				, created
				, user_id
				, resolution_criteria
				)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING id`

	err := r.db.QueryRowContext(ctx, query, f.Question, f.Category, f.CreatedAt, f.ResolutionCriteria).Scan(&f.ID)
	return err
}

func (r *ForecastRepository) ListOpenForecasts(ctx context.Context) ([]*models.Forecast, error) {
	query := `SELECT 
				id
				, question
				, category
				, created
				, user_id
				, resolution_criteria
				, resolution
				, resolved
				, comment 
				FROM forecast_v2
				WHERE resolved is null`

	return r.queryForecasts(ctx, query)
}

func (r *ForecastRepository) ListResolvedForecasts(ctx context.Context) ([]*models.Forecast, error) {
	query := `SELECT 
				id
				, question
				, category
				, created
				, user_id
				, resolution_criteria
				, resolution
				, resolved
				, comment 
				FROM forecast_v2
				WHERE resolved is not null`

	return r.queryForecasts(ctx, query)
}

func (r *ForecastRepository) ListOpenForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	query := `SELECT 
				id
				, question
				, category
				, created
				, user_id
				, resolution_criteria
				, resolution
				, resolved
				, comment 
				FROM forecast_v2
				WHERE resolved is null
				AND lower(category) like $1`
	categoryPattern := "%" + category + "%"

	return r.queryForecasts(ctx, query, categoryPattern)
}

func (r *ForecastRepository) ListResolvedForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	query := `SELECT
				id
				, question
				, category
				, created
				, user_id
				, resolution_criteria
				, resolution
				, resolved
				, comment 
				FROM forecast_v2
				WHERE resolved is not null
				AND lower(category) like $1`
	categoryPattern := "%" + category + "%"

	return r.queryForecasts(ctx, query, categoryPattern)
}

func (r *ForecastRepository) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE forecast_v2 SET 
				question = $1
				, category = $2
				, resolution_criteria = $3
				, resolution = $4
				, resolved = $5
				, comment = $9
			 WHERE id = $10`

	_, err = r.db.ExecContext(ctx, query,
		f.Question,
		f.Category,
		f.ResolutionCriteria,
		f.Resolution,
		f.ResolvedAt,
		f.ResolutionComment,
		f.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// user_id filtered methods
func (r *ForecastRepository) DeleteForecast(ctx context.Context, id int64, user_id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	//if this delete fails, due to user_id not owning forecast, return and do not delete forecast points.
	//also checking this in service, but this adds redundancy since I delete all user forecast points.
	queryForecasts := `DELETE FROM forecast_v2 WHERE id = $1 and user_id = $2`
	_, err = tx.ExecContext(ctx, queryForecasts, id, user_id)
	if err != nil {
		return err
	}

	queryForecastPoints := `DELETE FROM forecast_points WHERE forecast_id = $1`
	_, err = tx.ExecContext(ctx, queryForecastPoints, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Helper function to query forecasts
func (r *ForecastRepository) queryForecasts(ctx context.Context, query string, args ...interface{}) ([]*models.Forecast, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecasts []*models.Forecast

	for rows.Next() {
		var f models.Forecast
		err := rows.Scan(
			&f.ID,
			&f.Question,
			&f.Category,
			&f.CreatedAt,
			&f.UserID,
			&f.ResolutionCriteria,
			&f.Resolution,
			&f.ResolvedAt,
			&f.ResolutionComment)
		if err != nil {
			return nil, err
		}
		forecasts = append(forecasts, &f)
	}
	return forecasts, rows.Err()
}
