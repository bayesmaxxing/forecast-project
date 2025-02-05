package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"database/sql"

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

func (r *ScoreRepository) CreateScore(ctx context.Context, score *models.Scores) error {
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

func (r *ScoreRepository) GetScoresByUserID(ctx context.Context, userID int64) ([]models.Scores, error) {
	query := `SELECT 
				id, brier_score, log2_score, logn_score, user_id, forecast_id, created
			  FROM scores 
			  WHERE user_id = $1
			  ORDER BY created DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
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

func (r *ScoreRepository) DeleteScore(ctx context.Context, scoreID int64) error {
	query := `DELETE FROM scores WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, scoreID)
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
func (r *ScoreRepository) GetOverallScores(ctx context.Context) (*models.OverallScores, error) {
	query := `WITH user_scores AS (
				SELECT 
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
	query := `WITH category_scores AS (
				SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
				COUNT(DISTINCT user_id) as total_users,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecast_v2 f
			  ON s.forecast_id = f.id
			  WHERE f.category = $1`

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

func (r *ScoreRepository) GetUserScores(ctx context.Context, userID int64) (*models.UserScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores
			  WHERE user_id = $1`

	var userScores models.UserScores
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&userScores.BrierScore,
		&userScores.Log2Score,
		&userScores.LogNScore,
		&userScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}

	userScores.UserID = userID

	return &userScores, nil
}

func (r *ScoreRepository) GetUserCategoryScores(ctx context.Context, userID int64, category string) (*models.UserCategoryScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn,
				COUNT(DISTINCT forecast_id) as total_forecasts
			  FROM scores s
			  JOIN forecast_v2 f
			  ON s.forecast_id = f.id
			  WHERE s.user_id = $1 AND f.category = $2`

	var userCategoryScores models.UserCategoryScores
	err := r.db.QueryRowContext(ctx, query, userID, category).Scan(
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
