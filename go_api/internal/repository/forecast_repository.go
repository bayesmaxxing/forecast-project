package repository

import (
	"context"
	"go_api/internal/database"
	"go_api/internal/models"
	"strings"

	_ "github.com/jackc/pgx/v5"
)

type ForecastRepository struct {
	db *database.DB
}

func NewForecastRepository(db *database.DB) *ForecastRepository {
	return &ForecastRepository{db: db}
}

func (r *ForecastRepository) GetForecastByID(ctx context.Context, id int64) (*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria
				resolution, resolved, brier_score, log2_score, logn_score, comment
				FROM forecast_v2 WHERE id = $1`

	var f models.Forecast
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID, &f.Question, &f.Category, &f.CreatedAt, &f.ResolutionCriteria,
		&f.Resolution, &f.ResolvedAt, &f.BrierScore, &f.Log2Score,
		&f.Log2Score, &f.LogNScore, &f.ResolutionComment,
	)

	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *ForecastRepository) CreateForecast(ctx context.Context, f *models.Forecast) error {
	query := `INSERT INTO forecast_v2 (question, category, created, resolution_criteria, 
				resolution, resolved, brier_score, log2_score, logn_score, comment)
				VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, f.Question, f.Category, f.CreatedAt, f.ResolutionCriteria).Scan(&f.ID)
	return err
}

func (r *ForecastRepository) UpdateForecast(ctx context.Context, f *models.Forecast) error {
	query := `UPDATE forecast_v2 SET question = $1, category = $2, resolution_criteria = $3, resolution = $4,
			 resolved = $5, brier_score = $6, log2_score = $7, logn_score = $8, comment = $9
			 WHERE id = $10`

	_, err := r.db.ExecContext(ctx, query, f.Question, f.Category, f.CreatedAt, f.ResolutionCriteria,
		f.Resolution, f.ResolvedAt, f.BrierScore, f.Log2Score, f.LogNScore, f.ResolutionComment, f.ID)
	return err
}

func (r *ForecastRepository) DeleteForecast(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	queryForecastPoints := `DELETE FROM forecast_points WHERE forecast_id = $1`
	_, err = tx.ExecContext(ctx, queryForecastPoints, id)
	if err != nil {
		return err
	}

	queryForecasts := `DELETE FROM forecast_v2 WHERE id = $1`
	_, err = tx.ExecContext(ctx, queryForecasts, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ForecastRepository) ListOpenForecasts(ctx context.Context) ([]*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria, resolution, resolved, brier_score,
				log2_score, logn_score, comment 
				FROM forecast_v2
				WHERE resolved is null`

	return r.queryForecasts(ctx, query)
}

func (r *ForecastRepository) ListResolvedForecasts(ctx context.Context) ([]*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria, resolution, resolved, brier_score,
				log2_score, logn_score, comment 
				FROM forecast_v2
				WHERE resolved is not null`

	return r.queryForecasts(ctx, query)
}

func (r *ForecastRepository) ListOpenForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria, resolution, resolved, brier_score,
				log2_score, logn_score, comment 
				FROM forecast_v2
				WHERE resolved is null
				AND category like (%$1%)`

	return r.queryForecasts(ctx, query, category)
}

func (r *ForecastRepository) ListResolvedForecastsWithCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria, resolution, resolved, brier_score,
				log2_score, logn_score, comment 
				FROM forecast_v2
				WHERE resolved is not null
				AND category like (%$1%)`

	return r.queryForecasts(ctx, query, category)
}

func (r *ForecastRepository) ListResolvedWithScoresAndCategory(ctx context.Context, category string) ([]*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria, resolution, resolved, brier_score,
				log2_score, logn_score, comment 
				FROM forecast_v2
				WHERE resolved is not null
				AND brier_score is not null
				AND category like (%$1%)`

	return r.queryForecasts(ctx, query, category)
}

func (r *ForecastRepository) ListResolvedWithScores(ctx context.Context) ([]*models.Forecast, error) {
	query := `SELECT id, question, category, created, resolution_criteria, resolution, resolved, brier_score,
				log2_score, logn_score, comment 
				FROM forecast_v2
				WHERE resolved is not null
				AND brier_score is not null`

	return r.queryForecasts(ctx, query)
}

func (r *ForecastRepository) queryForecasts(ctx context.Context, query string, args ...interface{}) ([]*models.Forecast, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var forecasts []*models.Forecast

	for rows.Next() {
		var f models.Forecast
		var err error
		if strings.Contains(query, "resolved is null") {
			err = rows.Scan(&f.ID, &f.Question, &f.Category, &f.CreatedAt, &f.ResolutionCriteria)
		} else {
			err = rows.Scan(&f.ID, &f.Question, &f.Category, &f.CreatedAt, &f.ResolutionCriteria, &f.Resolution,
				&f.ResolvedAt, &f.BrierScore, &f.Log2Score, &f.LogNScore, &f.ResolutionComment)
		}
		if err != nil {
			return nil, err
		}
		forecasts = append(forecasts, &f)
	}
	return forecasts, rows.Err()
}
