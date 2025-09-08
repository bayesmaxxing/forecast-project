package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
)

// ScoreRepository defines the interface for score data operations
type ScoreRepository interface {
	GetScoreByForecastID(ctx context.Context, forecastID int64) ([]models.Scores, error)
	GetScoreByForecastAndUser(ctx context.Context, forecastID int64, userID int64) (*models.Scores, error)
	GetAverageScoreByForecastID(ctx context.Context, forecastID int64) (*models.ScoreMetrics, error)
	CreateScore(ctx context.Context, score *models.Scores) error
	GetScoresByUserID(ctx context.Context, userID int64) ([]models.Scores, error)
	GetAllScores(ctx context.Context) ([]models.Scores, error)
	GetAverageScores(ctx context.Context) ([]models.Scores, error)
	UpdateScore(ctx context.Context, score *models.Scores) error
	DeleteScore(ctx context.Context, scoreID int64) error
	GetOverallScores(ctx context.Context) (*models.OverallScores, error)
	GetCategoryScores(ctx context.Context, category string) (*models.CategoryScores, error)
	GetCategoryScoresByUsers(ctx context.Context, category string) ([]models.UserCategoryScores, error)
	GetOverallScoresByUsers(ctx context.Context) ([]models.UserScores, error)
	GetUserCategoryScores(ctx context.Context, userID int64, category string) (*models.UserCategoryScores, error)
	GetUserOverallScores(ctx context.Context, userID int64) (*models.UserScores, error)
}

// PostgresScoreRepository implements the ScoreRepository interface
type PostgresScoreRepository struct {
	db *database.DB
}

// NewScoreRepository creates a new PostgresScoreRepository instance
func NewScoreRepository(db *database.DB) ScoreRepository {
	return &PostgresScoreRepository{db: db}
}

// Full schema operations
func (r *PostgresScoreRepository) GetScoreByForecastID(ctx context.Context, forecast_id int64) ([]models.Scores, error) {
	query := `SELECT 
					id 
					, brier_score
					, log2_score
					, logn_score
					, user_id
					, forecast_id
					, created
					FROM scores 
					WHERE forecast_id = $1`

	rows, err := r.db.QueryContext(ctx, query, forecast_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Scores
	for rows.Next() {
		var s models.Scores
		if err := rows.Scan(
			&s.ID,
			&s.BrierScore,
			&s.Log2Score,
			&s.LogNScore,
			&s.UserID,
			&s.ForecastID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, rows.Err()
}

func (r *PostgresScoreRepository) GetScoreByForecastAndUser(ctx context.Context, forecast_id int64, user_id int64) (*models.Scores, error) {
	query := `SELECT 
					id 
					, brier_score
					, log2_score
					, logn_score
					, user_id
					, forecast_id
					, created
					FROM scores 
					WHERE forecast_id = $1 AND user_id = $2`

	var score models.Scores
	err := r.db.QueryRowContext(ctx, query, forecast_id, user_id).Scan(
		&score.ID,
		&score.BrierScore,
		&score.Log2Score,
		&score.LogNScore,
		&score.UserID,
		&score.ForecastID,
		&score.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *PostgresScoreRepository) CreateScore(ctx context.Context, score *models.Scores) error {
	score.CreatedAt = time.Now()

	query := `INSERT INTO scores (brier_score, log2_score, logn_score, user_id, forecast_id, created)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		score.BrierScore,
		score.Log2Score,
		score.LogNScore,
		score.UserID,
		score.ForecastID,
		score.CreatedAt).Scan(&score.ID)
}

func (r *PostgresScoreRepository) GetScoresByUserID(ctx context.Context, user_id int64) ([]models.Scores, error) {
	query := `SELECT 
				id, brier_score, log2_score, logn_score, user_id, forecast_id, created
			  FROM scores 
			  WHERE user_id = $1
			  ORDER BY created DESC`

	rows, err := r.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Scores
	for rows.Next() {
		var s models.Scores
		if err := rows.Scan(
			&s.ID,
			&s.BrierScore,
			&s.Log2Score,
			&s.LogNScore,
			&s.UserID,
			&s.ForecastID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, rows.Err()
}

func (r *PostgresScoreRepository) GetAllScores(ctx context.Context) ([]models.Scores, error) {
	query := `SELECT 
				id, brier_score, log2_score, logn_score, user_id, forecast_id, created
			  FROM scores 
			  ORDER BY created DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Scores
	for rows.Next() {
		var s models.Scores
		if err := rows.Scan(
			&s.ID,
			&s.BrierScore,
			&s.Log2Score,
			&s.LogNScore,
			&s.UserID,
			&s.ForecastID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, rows.Err()
}

func (r *PostgresScoreRepository) GetAverageScores(ctx context.Context) ([]models.Scores, error) {
	query := `SELECT 0 as id, coalesce(AVG(brier_score), 0) as brier_score, coalesce(AVG(log2_score), 0) as log2_score, coalesce(AVG(logn_score), 0) as logn_score, 0 as user_id, forecast_id, max(created) as created
			  FROM scores
			  GROUP BY forecast_id, user_id, id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Scores
	for rows.Next() {
		var s models.Scores
		if err := rows.Scan(
			&s.ID,
			&s.BrierScore,
			&s.Log2Score,
			&s.LogNScore,
			&s.UserID,
			&s.ForecastID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, rows.Err()
}

func (r *PostgresScoreRepository) UpdateScore(ctx context.Context, score *models.Scores) error {
	query := `UPDATE scores 
			  SET brier_score = $1, log2_score = $2, logn_score = $3
			  WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		score.BrierScore,
		score.Log2Score,
		score.LogNScore,
		score.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresScoreRepository) DeleteScore(ctx context.Context, score_id int64) error {
	query := `DELETE FROM scores WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, score_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Aggregate score operations
// all users
func (r *PostgresScoreRepository) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	query := `SELECT 
				coalesce(AVG(brier_score), 0) as avg_brier,
				coalesce(AVG(log2_score), 0) as avg_log2,
				coalesce(AVG(logn_score), 0) as avg_logn,
				COUNT(DISTINCT user_id) as total_users,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores`

	var overallScores models.OverallScores
	err := r.db.QueryRowContext(ctx, query).Scan(
		&overallScores.BrierScore,
		&overallScores.Log2Score,
		&overallScores.LogNScore,
		&overallScores.TotalUsers,
		&overallScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}

	return &overallScores, nil
}

func (r *PostgresScoreRepository) GetCategoryScores(ctx context.Context, category string) (*models.CategoryScores, error) {
	query := `SELECT 
				coalesce(AVG(s.brier_score), 0) as avg_brier,
				coalesce(AVG(s.log2_score), 0) as avg_log2,
				coalesce(AVG(s.logn_score), 0) as avg_logn,
				COUNT(DISTINCT s.user_id) as total_users,
				COUNT(DISTINCT s.forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecasts f
			  ON s.forecast_id = f.id
			  WHERE lower(f.category) LIKE lower($1)`

	categoryPattern := "%" + category + "%"

	var categoryScores models.CategoryScores
	err := r.db.QueryRowContext(ctx, query, categoryPattern).Scan(
		&categoryScores.BrierScore,
		&categoryScores.Log2Score,
		&categoryScores.LogNScore,
		&categoryScores.TotalUsers,
		&categoryScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}
	categoryScores.Category = category

	return &categoryScores, nil
}

func (r *PostgresScoreRepository) GetCategoryScoresByUsers(ctx context.Context, category string) ([]models.UserCategoryScores, error) {
	query := `SELECT 
				coalesce(AVG(s.brier_score), 0) as avg_brier,
				coalesce(AVG(s.log2_score), 0) as avg_log2,
				coalesce(AVG(s.logn_score), 0) as avg_logn,
				s.user_id,
				COUNT(DISTINCT s.forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecasts f
			  ON s.forecast_id = f.id
			  WHERE lower(f.category) LIKE lower($1)
			  GROUP BY s.user_id`

	categoryPattern := "%" + category + "%"

	rows, err := r.db.QueryContext(ctx, query, categoryPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCategoryScores []models.UserCategoryScores
	for rows.Next() {
		var u models.UserCategoryScores
		if err := rows.Scan(
			&u.BrierScore,
			&u.Log2Score,
			&u.LogNScore,
			&u.UserID,
			&u.TotalForecasts,
		); err != nil {
			return nil, err
		}
		u.Category = category
		userCategoryScores = append(userCategoryScores, u)
	}
	return userCategoryScores, rows.Err()
}

func (r *PostgresScoreRepository) GetOverallScoresByUsers(ctx context.Context) ([]models.UserScores, error) {
	query := `SELECT 
				coalesce(AVG(brier_score), 0) as avg_brier,
				coalesce(AVG(log2_score), 0) as avg_log2,
				coalesce(AVG(logn_score), 0) as avg_logn,
				user_id,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores
			  GROUP BY user_id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userScores []models.UserScores
	for rows.Next() {
		var u models.UserScores
		if err := rows.Scan(
			&u.BrierScore,
			&u.Log2Score,
			&u.LogNScore,
			&u.UserID,
			&u.TotalForecasts,
		); err != nil {
			return nil, err
		}
		userScores = append(userScores, u)
	}
	return userScores, rows.Err()
}

// user-specific
func (r *PostgresScoreRepository) GetUserCategoryScores(ctx context.Context, userID int64, category string) (*models.UserCategoryScores, error) {
	query := `SELECT 
				coalesce(AVG(s.brier_score), 0) as avg_brier,
				coalesce(AVG(s.log2_score), 0) as avg_log2,
				coalesce(AVG(s.logn_score), 0) as avg_logn,
				COUNT(DISTINCT s.forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecasts f
			  ON s.forecast_id = f.id
			  WHERE s.user_id = $1 AND lower(f.category) LIKE lower($2)`

	categoryPattern := "%" + category + "%"

	var userCategoryScores models.UserCategoryScores
	err := r.db.QueryRowContext(ctx, query, userID, categoryPattern).Scan(
		&userCategoryScores.BrierScore,
		&userCategoryScores.Log2Score,
		&userCategoryScores.LogNScore,
		&userCategoryScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}
	userCategoryScores.Category = category
	userCategoryScores.UserID = userID

	return &userCategoryScores, nil
}

func (r *PostgresScoreRepository) GetUserOverallScores(ctx context.Context, user_id int64) (*models.UserScores, error) {
	query := `SELECT 
				coalesce(AVG(s.brier_score), 0) as avg_brier,
				coalesce(AVG(s.log2_score), 0) as avg_log2,
				coalesce(AVG(s.logn_score), 0) as avg_logn,
				COUNT(DISTINCT s.forecast_id) as total_forecasts
			  FROM scores s
			  WHERE s.user_id = $1`

	var userScores models.UserScores
	err := r.db.QueryRowContext(ctx, query, user_id).Scan(
		&userScores.BrierScore,
		&userScores.Log2Score,
		&userScores.LogNScore,
		&userScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}

	return &userScores, nil
}

func (r *PostgresScoreRepository) GetAverageScoreByForecastID(ctx context.Context, forecast_id int64) (*models.ScoreMetrics, error) {
	query := `SELECT 
				coalesce(AVG(brier_score), 0) as avg_brier,
				coalesce(AVG(log2_score), 0) as avg_log2,
				coalesce(AVG(logn_score), 0) as avg_logn
			  FROM scores
			  WHERE forecast_id = $1`

	var scoreMetrics models.ScoreMetrics
	err := r.db.QueryRowContext(ctx, query, forecast_id).Scan(
		&scoreMetrics.BrierScore,
		&scoreMetrics.Log2Score,
		&scoreMetrics.LogNScore,
	)
	if err != nil {
		return nil, err
	}

	return &scoreMetrics, nil
}
