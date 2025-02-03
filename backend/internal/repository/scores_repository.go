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

func (r *ScoreRepository) GetScoreByForecastID(ctx context.Context, forecast_id int64) (*models.Scores, error) {
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

	var s models.Scores
	err := r.db.QueryRowContext(ctx, query, forecast_id).Scan(
		&s.ID,
		&s.BrierScore,
		&s.Log2Score,
		&s.LogNScore,
		&s.UserID,
		&s.ForecastID,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
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

func (r *ScoreRepository) GetAverageScoresAllUsers(ctx context.Context) ([]models.AverageScores, error) {
	query := `SELECT user_id,
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn
			  FROM scores 
			  GROUP BY user_id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var avgScores []models.AverageScores
	for rows.Next() {
		var score models.AverageScores
		if err := rows.Scan(
			&score.UserID,
			&score.AvgBrierScore,
			&score.AvgLog2Score,
			&score.AvgLogNScore,
		); err != nil {
			return nil, err
		}
		avgScores = append(avgScores, score)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return avgScores, nil
}

func (r *ScoreRepository) GetAverageScoresByUserID(ctx context.Context, userID int64) (*models.AverageScores, error) {
	query := `SELECT 
				AVG(brier_score) as avg_brier,
				AVG(log2_score) as avg_log2,
				AVG(logn_score) as avg_logn
			  FROM scores 
			  WHERE user_id = $1`

	var avgScores models.AverageScores
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&avgScores.UserID,
		&avgScores.AvgBrierScore,
		&avgScores.AvgLog2Score,
		&avgScores.AvgLogNScore,
	)
	if err != nil {
		return nil, err
	}

	return &avgScores, nil
}
