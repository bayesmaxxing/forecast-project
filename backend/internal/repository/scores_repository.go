package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
)

type ScoreRepository struct {
	db *database.DB
}

func NewScoreRepository(db *database.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

// Full schema operations
func (r *ScoreRepository) GetScoreByForecastID(ctx context.Context, forecast_id int64) ([]models.Scores, error) {
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

func (r *ScoreRepository) GetScoreByForecastAndUser(ctx context.Context, forecast_id int64, user_id int64) (*models.Scores, error) {
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

func (r *ScoreRepository) CreateScore(ctx context.Context, score *models.Scores) error {
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

func (r *ScoreRepository) GetScoresByUserID(ctx context.Context, user_id int64) ([]models.Scores, error) {
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

func (r *ScoreRepository) GetAllScores(ctx context.Context) ([]models.Scores, error) {
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

func (r *ScoreRepository) UpdateScore(ctx context.Context, score *models.Scores) error {
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

func (r *ScoreRepository) DeleteScore(ctx context.Context, score_id int64) error {
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
func (r *ScoreRepository) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
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

func (r *ScoreRepository) GetCategoryScores(ctx context.Context, category string) (*models.CategoryScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
				COUNT(DISTINCT user_id) as total_users,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecast_v2 f
			  ON s.forecast_id = f.id
			  WHERE f.category LIKE ($1)`

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

	return &categoryScores, nil
}

func (r *ScoreRepository) GetCategoryScoresByUsers(ctx context.Context, category string) ([]models.UserCategoryScores, error) {
	query := `SELECT 
				AVG(s.brier_score) as avg_brier,
				AVG(s.log2_score) as avg_log2,
				AVG(s.logn_score) as avg_logn,
				s.user_id,
				COUNT(DISTINCT s.forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecast_v2 f
			  ON s.forecast_id = f.id
			  WHERE f.category LIKE ($1)
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
			&u.Category,
			&u.TotalForecasts,
		); err != nil {
			return nil, err
		}
		userCategoryScores = append(userCategoryScores, u)
	}
	return userCategoryScores, rows.Err()
}

func (r *ScoreRepository) GetOverallScoresByUsers(ctx context.Context) ([]models.UserScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
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
func (r *ScoreRepository) GetUserCategoryScores(ctx context.Context, userID int64, category string) (*models.UserCategoryScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecast_v2 f
			  ON s.forecast_id = f.id
			  WHERE s.user_id = $1 AND f.category LIKE ($2)`

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

	return &userCategoryScores, nil
}

func (r *ScoreRepository) GetUserOverallScores(ctx context.Context, user_id int64) (*models.UserScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecast_v2 f
			  ON s.forecast_id = f.id
			  WHERE s.user_id = $1`

	var userScores models.UserScores
	err := r.db.QueryRowContext(ctx, query, user_id).Scan(
		&userScores.BrierScore,
		&userScores.Log2Score,
		&userScores.LogNScore,
		&userScores.UserID,
		&userScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}

	return &userScores, nil
}
