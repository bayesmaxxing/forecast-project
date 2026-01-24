package repository

import (
	"backend/internal/database"
	"backend/internal/logger"
	"backend/internal/models"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// CalibrationRepository defines the interface for calibration data operations
type CalibrationRepository interface {
	GetCalibrationData(ctx context.Context, filters models.CalibrationFilters) (*models.CalibrationData, error)
	GetCalibrationDataByUsers(ctx context.Context, filters models.CalibrationFilters) ([]models.UserCalibrationData, error)
}

// PostgresCalibrationRepository implements the CalibrationRepository interface
type PostgresCalibrationRepository struct {
	db *database.DB
}

// NewCalibrationRepository creates a new PostgresCalibrationRepository instance
func NewCalibrationRepository(db *database.DB) CalibrationRepository {
	return &PostgresCalibrationRepository{db: db}
}

func buildCalibrationBaseQuery(filters models.CalibrationFilters, groupByUser bool) (string, []any) {
	args := []any{}
	argsCounter := 1

	whereConditions := []string{"f.resolution IN ('0', '1')"}

	if filters.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.user_id = $%d", argsCounter))
		args = append(args, *filters.UserID)
		argsCounter++
	}
	if filters.Category != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("lower(f.category) LIKE $%d", argsCounter))
		args = append(args, "%"+*filters.Category+"%")
		argsCounter++
	}
	if filters.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created >= $%d", argsCounter))
		args = append(args, *filters.StartDate)
		argsCounter++
	}
	if filters.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created <= $%d", argsCounter))
		args = append(args, *filters.EndDate)
		argsCounter++
	}

	userIDSelect := ""
	userIDGroupBy := ""
	if groupByUser {
		userIDSelect = "p.user_id,"
		userIDGroupBy = "p.user_id,"
	}

	query := fmt.Sprintf(`
WITH point_outcomes AS (
    SELECT
        p.point_forecast,
        p.user_id,
        p.forecast_id,
        CASE WHEN f.resolution = '1' THEN 1 ELSE 0 END as outcome
    FROM points p
    INNER JOIN forecasts f ON p.forecast_id = f.id
    WHERE %s
)
SELECT
    %s
    FLOOR(point_forecast * 10) / 10 as bucket_start,
    COUNT(*) as prediction_count,
    AVG(point_forecast) as avg_prediction,
    AVG(outcome::float) as actual_rate,
    COUNT(DISTINCT forecast_id) as forecast_count
FROM point_outcomes
GROUP BY %s FLOOR(point_forecast * 10) / 10
ORDER BY %s bucket_start
`, strings.Join(whereConditions, " AND "), userIDSelect, userIDGroupBy, userIDSelect)

	return query, args
}

func (r *PostgresCalibrationRepository) GetCalibrationData(ctx context.Context, filters models.CalibrationFilters) (*models.CalibrationData, error) {
	log := logger.FromContext(ctx)

	query, args := buildCalibrationBaseQuery(filters, false)

	start := time.Now()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute calibration query", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("executed calibration query", slog.Duration("duration", time.Since(start)))
	defer rows.Close()

	var buckets []models.CalibrationBucket
	var totalPredictions, totalForecasts int

	for rows.Next() {
		var bucket models.CalibrationBucket
		var forecastCount int
		if err := rows.Scan(
			&bucket.BucketStart,
			&bucket.PredictionCount,
			&bucket.AvgPrediction,
			&bucket.ActualRate,
			&forecastCount,
		); err != nil {
			log.Error("failed to scan calibration row", slog.String("error", err.Error()))
			return nil, err
		}
		bucket.BucketEnd = bucket.BucketStart + 0.1
		buckets = append(buckets, bucket)
		totalPredictions += bucket.PredictionCount
		totalForecasts += forecastCount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Info("calibration query results", slog.Int("bucket_count", len(buckets)), slog.Int("total_predictions", totalPredictions))
	return &models.CalibrationData{
		Buckets:          buckets,
		TotalPredictions: totalPredictions,
		TotalForecasts:   totalForecasts,
	}, nil
}

func (r *PostgresCalibrationRepository) GetCalibrationDataByUsers(ctx context.Context, filters models.CalibrationFilters) ([]models.UserCalibrationData, error) {
	log := logger.FromContext(ctx)

	query, args := buildCalibrationBaseQuery(filters, true)

	start := time.Now()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute calibration by users query", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("executed calibration by users query", slog.Duration("duration", time.Since(start)))
	defer rows.Close()

	// Map to collect buckets per user
	userDataMap := make(map[int64]*models.UserCalibrationData)

	for rows.Next() {
		var userID int64
		var bucket models.CalibrationBucket
		var forecastCount int
		if err := rows.Scan(
			&userID,
			&bucket.BucketStart,
			&bucket.PredictionCount,
			&bucket.AvgPrediction,
			&bucket.ActualRate,
			&forecastCount,
		); err != nil {
			log.Error("failed to scan calibration by users row", slog.String("error", err.Error()))
			return nil, err
		}
		bucket.BucketEnd = bucket.BucketStart + 0.1

		if _, exists := userDataMap[userID]; !exists {
			userDataMap[userID] = &models.UserCalibrationData{
				UserID:  userID,
				Buckets: []models.CalibrationBucket{},
			}
		}
		userData := userDataMap[userID]
		userData.Buckets = append(userData.Buckets, bucket)
		userData.TotalPredictions += bucket.PredictionCount
		userData.TotalForecasts += forecastCount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	result := make([]models.UserCalibrationData, 0, len(userDataMap))
	for _, userData := range userDataMap {
		result = append(result, *userData)
	}

	log.Info("calibration by users query results", slog.Int("user_count", len(result)))
	return result, nil
}
