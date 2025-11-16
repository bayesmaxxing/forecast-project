package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5"
)

// ForecastPointRepository defines the interface for forecast point data operations
type ForecastPointRepository interface {
	GetForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error)
	CreateForecastPoint(ctx context.Context, fp *models.ForecastPoint) error
}

// PostgresForecastPointRepository implements the ForecastPointRepository interface
type PostgresForecastPointRepository struct {
	db *database.DB
}

// NewForecastPointRepository creates a new PostgresForecastPointRepository instance
func NewForecastPointRepository(db *database.DB) ForecastPointRepository {
	return &PostgresForecastPointRepository{db: db}
}

func buildForecastPointQuery(filters models.PointFilters) (string, error) {
	// start with distinct on forecast if requested
	// the fields to select are always the same
	distinctClause := []string{}
	if filters.DistinctOnForecast != nil && *filters.DistinctOnForecast {
		distinctClause = append(distinctClause, "distinct on (forecast_id)")
	} else {
		distinctClause = append(distinctClause, "")
	}

	selectFields := []string{
		"p.id",
		"p.forecast_id",
		"p.point_forecast",
		"p.reason",
		"p.created",
		"p.user_id",
		"u.username",
	}

	//so are from clauses
	fromClauses := []string{
		"points p",
		"left join users u",
		"on p.user_id = u.id",
	}

	// build where conditions
	argsCounter := 1
	whereConditions := []string{"1=1"}
	if filters.UserID != nil {
		whereConditions = append(whereConditions, "p.user_id = "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}
	if filters.ForecastID != nil {
		whereConditions = append(whereConditions, "p.forecast_id = "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}
	if filters.Date != nil {
		whereConditions = append(whereConditions, "p.created >= "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
		whereConditions = append(whereConditions, "p.created < "+fmt.Sprintf("$%d", argsCounter))
		argsCounter++
	}

	// finally order by statements
	orderBy := []string{}
	if filters.OrderByForecastID != nil && *filters.OrderByForecastID {
		orderBy = append(orderBy, "p.forecast_id")
	}

	if filters.CreatedDirection != nil {
		orderBy = append(orderBy, "p.created "+*filters.CreatedDirection)
	} else {
		orderBy = append(orderBy, "p.created DESC")
	}

	// build query
	query := fmt.Sprintf(
		`select
			%s
			%s
		from %s
		where %s
		order by %s`,
		strings.Join(distinctClause, " "),
		strings.Join(selectFields, ", "),
		strings.Join(fromClauses, "\n\t\t"),
		strings.Join(whereConditions, " and "),
		strings.Join(orderBy, ", "),
	)

	return query, nil
}

func (r *PostgresForecastPointRepository) GetForecastPoints(ctx context.Context, filters models.PointFilters) ([]*models.ForecastPoint, error) {
	query, err := buildForecastPointQuery(filters)
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
	if filters.Date != nil {
		args = append(args, *filters.Date)
		args = append(args, filters.Date.AddDate(0, 0, 1))
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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
