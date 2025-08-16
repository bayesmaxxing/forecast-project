package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"time"

	_ "github.com/jackc/pgx/v5"
)

// ForecastRepository defines the interface for forecast data operations
type ForecastRepository interface {
	GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error)
	CheckForecastOwnership(ctx context.Context, id int64, userID int64) (bool, error)
	CheckForecastStatus(ctx context.Context, id int64) (bool, error)
	CreateForecast(ctx context.Context, f *models.Forecast) error
	ListOpenForecasts(ctx context.Context) ([]*models.Forecast, error)
	ListResolvedForecasts(ctx context.Context) ([]*models.Forecast, error)
	ListOpenForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error)
	ListResolvedForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error)
	UpdateForecast(ctx context.Context, f *models.Forecast) error
	DeleteForecast(ctx context.Context, id int64, userID int64) error
	GetStaleAndNewForecasts(ctx context.Context, userID int64) ([]*models.Forecast, error)
}

// PostgresForecastRepository implements the ForecastRepository interface
type PostgresForecastRepository struct {
	db *database.DB
}

// NewForecastRepository creates a new PostgresForecastRepository instance
func NewForecastRepository(db *database.DB) ForecastRepository {
	return &PostgresForecastRepository{db: db}
}

// Methods without user_id filtering
func (r *PostgresForecastRepository) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
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
				FROM forecasts 
				WHERE id = $1`

	var f models.Forecast
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
	return &f, nil
}

func (r *PostgresForecastRepository) CheckForecastOwnership(ctx context.Context, id int64, user_id int64) (bool, error) {
	query := `SELECT user_id FROM forecasts WHERE id = $1`
	var forecastUserID int64
	err := r.db.QueryRowContext(ctx, query, id).Scan(&forecastUserID)
	if err != nil {
		return false, err
	}

	return forecastUserID == user_id, nil
}

func (r *PostgresForecastRepository) CheckForecastStatus(ctx context.Context, id int64) (bool, error) {
	query := `SELECT (resolved is null) FROM forecasts WHERE id = $1`
	var resolved bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&resolved)

	return resolved, err
}

func (r *PostgresForecastRepository) CreateForecast(ctx context.Context, f *models.Forecast) error {
	f.CreatedAt = time.Now()

	query := `INSERT INTO forecasts (
				question
				, category
				, created
				, user_id
				, resolution_criteria
				)
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING id`

	err := r.db.QueryRowContext(ctx, query, f.Question, f.Category, f.CreatedAt, f.UserID, f.ResolutionCriteria).Scan(&f.ID)
	return err
}

func (r *PostgresForecastRepository) ListOpenForecasts(ctx context.Context) ([]*models.Forecast, error) {
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
				FROM forecasts
				WHERE resolved is null`

	return r.queryForecasts(ctx, query)
}

func (r *PostgresForecastRepository) ListResolvedForecasts(ctx context.Context) ([]*models.Forecast, error) {
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
				FROM forecasts
				WHERE resolved is not null`

	return r.queryForecasts(ctx, query)
}

func (r *PostgresForecastRepository) ListOpenForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
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
				FROM forecasts
				WHERE resolved is null
				AND lower(category) like $1`
	categoryPattern := "%" + category + "%"

	return r.queryForecasts(ctx, query, categoryPattern)
}

func (r *PostgresForecastRepository) ListResolvedForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
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
				FROM forecasts
				WHERE resolved is not null
				AND lower(category) like $1`
	categoryPattern := "%" + category + "%"

	return r.queryForecasts(ctx, query, categoryPattern)
}

func (r *PostgresForecastRepository) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE forecasts SET 
				question = $1
				, category = $2
				, resolution_criteria = $3
				, resolution = $4
				, resolved = $5
				, comment = $6
			 WHERE id = $7`

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
func (r *PostgresForecastRepository) DeleteForecast(ctx context.Context, id int64, user_id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	//if this delete fails, due to user_id not owning forecast, return and do not delete forecast points.
	//also checking this in service, but this adds redundancy since I delete all user forecast points.
	queryForecasts := `DELETE FROM forecasts WHERE id = $1 and user_id = $2`
	_, err = tx.ExecContext(ctx, queryForecasts, id, user_id)
	if err != nil {
		return err
	}

	queryForecastPoints := `DELETE FROM points WHERE forecast_id = $1`
	_, err = tx.ExecContext(ctx, queryForecastPoints, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresForecastRepository) GetStaleAndNewForecasts(ctx context.Context, userID int64) ([]*models.Forecast, error) {
	query := `WITH latest_forecast_points as (
							 select forecast_id
							 , max(created) latest_created
							 , count(id) count_forecast_points
							 from points
							 where user_id = $1
							 group by forecast_id
							)

							SELECT
							f.id
							, f.question
							, f.category
							, f.created
							, f.user_id
							, f.resolution_criteria
							, f.resolution
							, f.resolved
							, f.comment
							FROM forecasts f
							LEFT JOIN latest_forecast_points lfp
							ON f.id = lfp.forecast_id
							WHERE f.resolved is null
							AND (lfp.forecast_id is null or lfp.latest_created < current_date - 30)
							limit 10`

	return r.queryForecasts(ctx, query, userID)
}

// Helper function to query forecasts
func (r *PostgresForecastRepository) queryForecasts(ctx context.Context, query string, args ...interface{}) ([]*models.Forecast, error) {
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
