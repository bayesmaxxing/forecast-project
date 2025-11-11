package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"time"

	_ "github.com/jackc/pgx/v5"
)

// ForecastPointRepository defines the interface for forecast point data operations
type ForecastPointRepository interface {
	GetAllForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error)
	GetForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error)
	GetForecastPointsByForecastIDAndUser(ctx context.Context, id int64, userID int64) ([]*models.ForecastPoint, error)
	GetOrderedForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error)
	CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error
	GetLatestForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error)
	GetLatestForecastPointsByUser(ctx context.Context, userID int64) ([]*models.ForecastPoint, error)
	GetForecastPointsByDate(ctx context.Context, userID int64, date *time.Time) ([]*models.ForecastPoint, error)
}

// PostgresForecastPointRepository implements the ForecastPointRepository interface
type PostgresForecastPointRepository struct {
	db *database.DB
}

// NewForecastPointRepository creates a new PostgresForecastPointRepository instance
func NewForecastPointRepository(db *database.DB) ForecastPointRepository {
	return &PostgresForecastPointRepository{db: db}
}

func (r *PostgresForecastPointRepository) GetAllForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.reason
				, p.created
				, p.user_id
				, u.username
				FROM points p
				left join users u on p.user_id = u.id
				ORDER BY p.created DESC`

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
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *PostgresForecastPointRepository) GetForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.reason
				, p.created
				, p.user_id
				, u.username
				FROM points p
				left join users u on p.user_id = u.id
				WHERE p.forecast_id = $1`

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
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *PostgresForecastPointRepository) GetForecastPointsByForecastIDAndUser(ctx context.Context, id int64, user_id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.reason
				, p.created
				, p.user_id
				, u.username
				FROM points p
				left join users u on p.user_id = u.id
				WHERE p.forecast_id = $1
				AND p.user_id = $2`

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
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *PostgresForecastPointRepository) CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error {
	fp.CreatedAt = time.Now()

	query := `INSERT INTO points (forecast_id
											, point_forecast
											, created
											, reason
											, user_id)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`

	err := r.db.QueryRowContext(ctx, query, fp.ForecastID, fp.PointForecast, fp.CreatedAt, fp.Reason, fp.UserID).Scan(&fp.ID)
	return err
}

func (r *PostgresForecastPointRepository) GetLatestForecastPoints(ctx context.Context) ([]*models.ForecastPoint, error) {
	query := `SELECT distinct on (forecast_id)
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.created
				, p.reason
				, p.user_id
				, u.username
				FROM points p
				left join users u on p.user_id = u.id
				ORDER BY p.forecast_id, p.created DESC;`

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
			&fp.CreatedAt,
			&fp.Reason,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *PostgresForecastPointRepository) GetLatestForecastPointsByUser(ctx context.Context, user_id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT distinct on (forecast_id)
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.created
				, p.reason
				, p.user_id
				, u.username
				FROM points p 
				left join users u on p.user_id = u.id
				WHERE p.user_id = $1
				ORDER BY p.forecast_id, p.created DESC;`

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
			&fp.CreatedAt,
			&fp.Reason,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *PostgresForecastPointRepository) GetOrderedForecastPointsByForecastID(ctx context.Context, id int64) ([]*models.ForecastPoint, error) {
	query := `SELECT 
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.reason
				, p.created
				, p.user_id
				, u.username
				FROM points p 
				left join users u on p.user_id = u.id
				WHERE p.forecast_id = $1
				ORDER BY p.id ASC`

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
			&fp.Reason,
			&fp.CreatedAt,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}

func (r *PostgresForecastPointRepository) GetForecastPointsByDate(ctx context.Context, user_id int64, date *time.Time) ([]*models.ForecastPoint, error) {
	// Use provided date or default to today
	targetDate := time.Now()
	if date != nil {
		targetDate = *date
	}
	// Truncate to start of day for date comparison
	targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	nextDay := targetDate.AddDate(0, 0, 1)

	query := `SELECT distinct on (forecast_id)
				p.id
				, p.forecast_id
				, p.point_forecast
				, p.created
				, p.reason
				, p.user_id
				, u.username
				FROM points p
				left join users u on p.user_id = u.id
				WHERE p.user_id = $1
				AND p.created >= $2
				AND p.created < $3
				ORDER BY p.forecast_id, p.created DESC;`

	rows, err := r.db.QueryContext(ctx, query, user_id, targetDate, nextDay)
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
			&fp.CreatedAt,
			&fp.Reason,
			&fp.UserID,
			&fp.UserName); err != nil {
			return nil, err
		}
		forecast_points = append(forecast_points, &fp)
	}
	return forecast_points, rows.Err()
}
