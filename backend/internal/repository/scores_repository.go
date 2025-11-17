package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5"
)

// ScoreRepository defines the interface for score data operations
type ScoreRepository interface {
	// get scores
	GetScores(ctx context.Context, filters models.ScoreFilters) ([]models.Scores, error)
	GetAverageScores(ctx context.Context) ([]models.Scores, error)

	// create, update, delete scores
	CreateScore(ctx context.Context, score *models.Scores) error
	UpdateScore(ctx context.Context, score *models.Scores) error
	DeleteScore(ctx context.Context, scoreID int64) error

	// aggregate scores
	GetAggregateScores(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error)
	GetAggregateScoresByUsers(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error)
}

// PostgresScoreRepository implements the ScoreRepository interface
type PostgresScoreRepository struct {
	db *database.DB
}

// NewScoreRepository creates a new PostgresScoreRepository instance
func NewScoreRepository(db *database.DB) ScoreRepository {
	return &PostgresScoreRepository{db: db}
}

func buildScoreQuery(filters models.ScoreFilters) (string, error) {
	selectFields := []string{
		"id",
		"brier_score",
		"log2_score",
		"logn_score",
		"brier_score_time_weighted",
		"log2_score_time_weighted",
		"logn_score_time_weighted",
		"user_id",
		"forecast_id",
		"created",
	}

	fromClause := "scores"

	whereConditions := []string{"1=1"}
	argsCounter := 1
	if filters.UserID != nil {
		whereConditions = append(whereConditions, "user_id = "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}
	if filters.ForecastID != nil {
		whereConditions = append(whereConditions, "forecast_id = "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}

	orderBy := "created DESC"

	query := fmt.Sprintf(
		`select %s from %s where %s order by %s`,
		strings.Join(selectFields, ", "),
		fromClause,
		strings.Join(whereConditions, " and "),
		orderBy,
	)
	fmt.Println(query)
	return query, nil
}

func (r *PostgresScoreRepository) GetScores(ctx context.Context, filters models.ScoreFilters) ([]models.Scores, error) {
	query, err := buildScoreQuery(filters)
	if err != nil {
		return nil, err
	}

	args := []any{}
	if filters.UserID != nil {
		args = append(args, *filters.UserID)
	}
	if filters.ForecastID != nil {
		args = append(args, *filters.ForecastID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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
			&s.BrierScoreTimeWeighted,
			&s.Log2ScoreTimeWeighted,
			&s.LogNScoreTimeWeighted,
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

func (r *PostgresScoreRepository) CreateScore(ctx context.Context, score *models.Scores) error {
	score.CreatedAt = time.Now()

	query := `INSERT INTO scores (brier_score
					, log2_score
					, logn_score
					, brier_score_time_weighted
					, log2_score_time_weighted
					, logn_score_time_weighted
					, user_id
					, forecast_id
					, created)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		score.BrierScore,
		score.Log2Score,
		score.LogNScore,
		score.BrierScoreTimeWeighted,
		score.Log2ScoreTimeWeighted,
		score.LogNScoreTimeWeighted,
		score.UserID,
		score.ForecastID,
		score.CreatedAt).Scan(&score.ID)
}

func (r *PostgresScoreRepository) GetAverageScores(ctx context.Context) ([]models.Scores, error) {
	query := `SELECT 0 as id
				, coalesce(AVG(brier_score), 0) as brier_score
				, coalesce(AVG(log2_score), 0) as log2_score
				, coalesce(AVG(logn_score), 0) as logn_score
				, coalesce(AVG(brier_score_time_weighted), 0) as brier_score_time_weighted
				, coalesce(AVG(log2_score_time_weighted), 0) as log2_score_time_weighted
				, coalesce(AVG(logn_score_time_weighted), 0) as logn_score_time_weighted
				, 0 as user_id
				, forecast_id
				, max(created) as created
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
			&s.BrierScoreTimeWeighted,
			&s.Log2ScoreTimeWeighted,
			&s.LogNScoreTimeWeighted,
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
			  SET brier_score = $1
			  , log2_score = $2
			  , logn_score = $3
			  , brier_score_time_weighted = $4
			  , log2_score_time_weighted = $5
			  , logn_score_time_weighted = $6
			  WHERE id = $7`

	result, err := r.db.ExecContext(ctx, query,
		score.BrierScore,
		score.Log2Score,
		score.LogNScore,
		score.BrierScoreTimeWeighted,
		score.Log2ScoreTimeWeighted,
		score.LogNScoreTimeWeighted,
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

func buildAggregateScoreQuery(filters models.ScoreFilters) (string, error) {
	// build select fields, start with the common ones
	selectFields := []string{
		"coalesce(AVG(s.brier_score), 0) as avg_brier",
		"coalesce(AVG(s.log2_score), 0) as avg_log2",
		"coalesce(AVG(s.logn_score), 0) as avg_logn",
		"coalesce(AVG(s.brier_score_time_weighted), 0) as avg_brier_time_weighted",
		"coalesce(AVG(s.log2_score_time_weighted), 0) as avg_log2_time_weighted",
		"coalesce(AVG(s.logn_score_time_weighted), 0) as avg_logn_time_weighted",
	}

	groupByClauses := []string{}
	if filters.GroupByUserID != nil && *filters.GroupByUserID {
		selectFields = append(selectFields, "s.user_id")
		selectFields = append(selectFields, "COUNT(DISTINCT s.forecast_id) as total_forecasts")
		groupByClauses = append(groupByClauses, "group by s.user_id")
	} else {
		selectFields = append(selectFields, "COUNT(DISTINCT s.user_id) as total_users")
		selectFields = append(selectFields, "COUNT(DISTINCT s.forecast_id) as total_forecasts")
	}

	fromClause := "scores s"

	joinClauses := []string{}
	whereConditions := []string{"1=1"}
	argsCounter := 1
	if filters.UserID != nil {
		whereConditions = append(whereConditions, "s.user_id = "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}
	if filters.ForecastID != nil {
		whereConditions = append(whereConditions, "s.forecast_id = "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}
	if filters.Category != nil {
		joinClauses = append(joinClauses, "left join forecasts f on s.forecast_id = f.id")
		whereConditions = append(whereConditions, "lower(f.category) like "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}

	query := fmt.Sprintf(
		`select %s from %s %s where %s %s`,
		strings.Join(selectFields, ", "),
		fromClause,
		strings.Join(joinClauses, "\n\t\t"),
		strings.Join(whereConditions, " and "),
		strings.Join(groupByClauses, "\n\t\t"),
	)
	return query, nil
}

func (r *PostgresScoreRepository) GetAggregateScores(ctx context.Context, filters models.ScoreFilters) (*models.OverallScores, error) {
	query, err := buildAggregateScoreQuery(filters)
	if err != nil {
		return nil, err
	}

	args := []any{}
	if filters.UserID != nil {
		args = append(args, *filters.UserID)
	}
	if filters.ForecastID != nil {
		args = append(args, *filters.ForecastID)
	}
	if filters.Category != nil {
		categoryPattern := "%" + *filters.Category + "%"
		args = append(args, categoryPattern)
	}
	if filters.GroupByUserID != nil && *filters.GroupByUserID {
		return nil, errors.New("group by user id is not supported")
	}

	var aggregateScores models.OverallScores
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&aggregateScores.BrierScore,
		&aggregateScores.Log2Score,
		&aggregateScores.LogNScore,
		&aggregateScores.BrierScoreTimeWeighted,
		&aggregateScores.Log2ScoreTimeWeighted,
		&aggregateScores.LogNScoreTimeWeighted,
		&aggregateScores.TotalUsers,
		&aggregateScores.TotalForecasts,
	)
	if err != nil {
		return nil, err
	}
	return &aggregateScores, nil
}

func (r *PostgresScoreRepository) GetAggregateScoresByUsers(ctx context.Context, filters models.ScoreFilters) ([]models.UserScores, error) {
	query, err := buildAggregateScoreQuery(filters)
	fmt.Println(query)
	if err != nil {
		return nil, err
	}

	args := []any{}
	if filters.UserID != nil {
		args = append(args, *filters.UserID)
	}
	if filters.Category != nil {
		categoryPattern := "%" + *filters.Category + "%"
		args = append(args, categoryPattern)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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
			&u.BrierScoreTimeWeighted,
			&u.Log2ScoreTimeWeighted,
			&u.LogNScoreTimeWeighted,
			&u.UserID,
			&u.TotalForecasts,
		); err != nil {
			return nil, err
		}
		userScores = append(userScores, u)
	}
	return userScores, rows.Err()
}
