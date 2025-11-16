package repository

import (
	"backend/internal/models"
	"strings"
	"testing"
	"time"
)

// Forecast Point Queries tests
func TestBuildForecastPointQueryAllFilters(t *testing.T) {

	userID := int64(2)
	forecastID := int64(3)
	date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	distinctOnForecast := true
	orderByForecastID := true
	createdDirection := "DESC"

	expectedQuery := `select
		distinct on (forecast_id)
			p.id,
			p.forecast_id,
			p.point_forecast,
			p.reason,
			p.created,
			p.user_id,
			u.username
		from points p
		left join users u on p.user_id = u.id
		where 1=1 
		and p.user_id = $1
		and p.forecast_id = $2
		and p.created >= $3
		and p.created < $4
		order by p.forecast_id, p.created DESC`

	filters := models.PointFilters{
		UserID:             &userID,
		ForecastID:         &forecastID,
		Date:               &date,
		DistinctOnForecast: &distinctOnForecast,
		OrderByForecastID:  &orderByForecastID,
		CreatedDirection:   &createdDirection,
	}

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	// Normalize both queries for comparison
	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

func TestBuildForecastPointQueryNoFilters(t *testing.T) {
	filters := models.PointFilters{}
	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	expectedQuery := `select
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	from points p
	left join users u on p.user_id = u.id
	where 1=1
	order by p.created DESC`

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

func TestBuildForecastPointQueryAscDirection(t *testing.T) {
	createdDirection := "ASC"
	orderByForecastID := false
	filters := models.PointFilters{
		CreatedDirection:  &createdDirection,
		OrderByForecastID: &orderByForecastID,
	}

	expectedQuery := `select
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	from points p
	left join users u on p.user_id = u.id
	where 1=1
	order by p.created ASC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

func TestBuildForecastPointForecastID(t *testing.T) {
	forecastID := int64(3)
	filters := models.PointFilters{
		ForecastID: &forecastID,
	}
	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	expectedQuery := `select
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	from points p
	left join users u on p.user_id = u.id
	where 1=1
	and p.forecast_id = $1
	order by p.created DESC`

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_GetForecastPointsByForecastID tests the query builder
// for the GetForecastPointsByForecastID method pattern
func TestBuildForecastPointQuery_GetForecastPointsByForecastID(t *testing.T) {
	forecastID := int64(5)
	filters := models.PointFilters{
		ForecastID: &forecastID,
	}

	expectedQuery := `SELECT 
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p
	left join users u on p.user_id = u.id
	WHERE 1=1 and p.forecast_id = $1
	ORDER BY p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_GetForecastPointsByForecastIDAndUser tests the query builder
// for the GetForecastPointsByForecastIDAndUser method pattern
func TestBuildForecastPointQuery_GetForecastPointsByForecastIDAndUser(t *testing.T) {
	forecastID := int64(5)
	userID := int64(10)
	filters := models.PointFilters{
		ForecastID: &forecastID,
		UserID:     &userID,
	}

	expectedQuery := `SELECT 
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p
	left join users u on p.user_id = u.id
	WHERE 1=1 
	and p.user_id = $1
	AND p.forecast_id = $2
	ORDER BY p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_GetLatestForecastPoints tests the query builder
// for the GetLatestForecastPoints method pattern
func TestBuildForecastPointQuery_GetLatestForecastPoints(t *testing.T) {
	distinctOnForecast := true
	orderByForecastID := true
	filters := models.PointFilters{
		DistinctOnForecast: &distinctOnForecast,
		OrderByForecastID:  &orderByForecastID,
	}

	expectedQuery := `SELECT distinct on (forecast_id)
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p
	left join users u on p.user_id = u.id
	WHERE 1=1
	ORDER BY p.forecast_id, p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_GetLatestForecastPointsByUser tests the query builder
// for the GetLatestForecastPointsByUser method pattern
func TestBuildForecastPointQuery_GetLatestForecastPointsByUser(t *testing.T) {
	userID := int64(7)
	distinctOnForecast := true
	orderByForecastID := true
	filters := models.PointFilters{
		UserID:             &userID,
		DistinctOnForecast: &distinctOnForecast,
		OrderByForecastID:  &orderByForecastID,
	}

	expectedQuery := `SELECT distinct on (forecast_id)
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p 
	left join users u on p.user_id = u.id
	WHERE 1=1 and p.user_id = $1
	ORDER BY p.forecast_id, p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_GetForecastPointsByDate tests the query builder
// for the GetForecastPointsByDate method pattern
func TestBuildForecastPointQuery_GetForecastPointsByDate(t *testing.T) {
	userID := int64(7)
	date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	distinctOnForecast := true
	orderByForecastID := true
	filters := models.PointFilters{
		UserID:             &userID,
		Date:               &date,
		DistinctOnForecast: &distinctOnForecast,
		OrderByForecastID:  &orderByForecastID,
	}

	expectedQuery := `SELECT distinct on (forecast_id)
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p
	left join users u on p.user_id = u.id
	WHERE 1=1
	AND p.user_id = $1
	AND p.created >= $2
	AND p.created < $3
	ORDER BY p.forecast_id, p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_WithUserIDOnly tests filtering by user only
func TestBuildForecastPointQuery_WithUserIDOnly(t *testing.T) {
	userID := int64(3)
	filters := models.PointFilters{
		UserID: &userID,
	}

	expectedQuery := `SELECT 
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p
	left join users u on p.user_id = u.id
	WHERE 1=1 and p.user_id = $1
	ORDER BY p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_WithDateOnly tests filtering by date only
func TestBuildForecastPointQuery_WithDateOnly(t *testing.T) {
	date := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	filters := models.PointFilters{
		Date: &date,
	}

	expectedQuery := `SELECT 
		p.id,
		p.forecast_id,
		p.point_forecast,
		p.reason,
		p.created,
		p.user_id,
		u.username
	FROM points p
	left join users u on p.user_id = u.id
	WHERE 1=1
	AND p.created >= $1
	AND p.created < $2
	ORDER BY p.created DESC`

	query, err := buildForecastPointQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast point query: %v", err)
	}

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// TestBuildForecastPointQuery_OrderingVariations tests different ordering options
func TestBuildForecastPointQuery_OrderingVariations(t *testing.T) {
	tests := []struct {
		name            string
		filters         models.PointFilters
		expectedOrderBy string
	}{
		{
			name:            "Default DESC ordering",
			filters:         models.PointFilters{},
			expectedOrderBy: "ORDER BY p.created DESC",
		},
		{
			name: "ASC ordering",
			filters: models.PointFilters{
				CreatedDirection: stringPtr("ASC"),
			},
			expectedOrderBy: "ORDER BY p.created ASC",
		},
		{
			name: "With forecast_id ordering",
			filters: models.PointFilters{
				OrderByForecastID: boolPtr(true),
			},
			expectedOrderBy: "ORDER BY p.forecast_id, p.created DESC",
		},
		{
			name: "With forecast_id and ASC",
			filters: models.PointFilters{
				OrderByForecastID: boolPtr(true),
				CreatedDirection:  stringPtr("ASC"),
			},
			expectedOrderBy: "ORDER BY p.forecast_id, p.created ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := buildForecastPointQuery(tt.filters)
			if err != nil {
				t.Fatalf("Error building forecast point query: %v", err)
			}

			normalizedQuery := normalizeSQL(query)
			if !strings.Contains(normalizedQuery, normalizeSQL(tt.expectedOrderBy)) {
				t.Errorf("Expected query to contain '%s', but got: %s",
					normalizeSQL(tt.expectedOrderBy), normalizedQuery)
			}
		})
	}
}

// Forecast queries tests
func TestBuildForecastQueryAllFilters(t *testing.T) {
	forecastID := int64(1)
	category := "finance"
	status := "open"
	filters := models.ForecastFilters{
		ForecastID: &forecastID,
		Status:     &status,
		Category:   &category,
	}

	query, err := buildForecastQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast query: %v", err)
	}

	expectedQuery := `select
		id,
		question,
		category,
		created,
		user_id,
		resolution_criteria,
		closing_date,
		resolution,
		resolved, 
		comment 
		from forecasts
		where 1=1 and id = $1 
		and resolved is null
		and lower(category) like $2`

	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

func TestBuildForecastQuery_WithClosedStatus(t *testing.T) {
	status := "closed"
	filters := models.ForecastFilters{
		Status: &status,
	}
	query, err := buildForecastQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast query: %v", err)
	}
	expectedQuery := `select
		id,
		question,
		category,
		created,
		user_id,
		resolution_criteria,
		closing_date,
		resolution,
		resolved, 
		comment 
		from forecasts
		where 1=1
		and current_date > closing_date`
	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

func TestBuildForecastQuery_WithCategoryAndResolvedStatus(t *testing.T) {
	category := "finance"
	status := "resolved"
	filters := models.ForecastFilters{
		Category: &category,
		Status:   &status,
	}
	query, err := buildForecastQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast query: %v", err)
	}
	expectedQuery := `select
		id,
		question,
		category,
		created,
		user_id,
		resolution_criteria,
		closing_date,
		resolution,
		resolved, 
		comment 
		from forecasts
		where 1=1
		and resolved is not null
		and lower(category) like $1`
	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

func TestBuildForecastQuery_NoFilters(t *testing.T) {
	filters := models.ForecastFilters{}
	query, err := buildForecastQuery(filters)
	if err != nil {
		t.Fatalf("Error building forecast query: %v", err)
	}
	expectedQuery := `select
		id,
		question,
		category,
		created,
		user_id,
		resolution_criteria,
		closing_date,
		resolution,
		resolved, 
		comment 
		from forecasts
		where 1=1`
	normalizedExpected := normalizeSQL(expectedQuery)
	normalizedActual := normalizeSQL(query)

	if normalizedActual != normalizedExpected {
		t.Errorf("Query mismatch:\nExpected: %s\nGot: %s", normalizedExpected, normalizedActual)
	}
}

// Helper functions for test data
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func normalizeSQL(sql string) string {
	// Replace multiple whitespace characters with a single space
	sql = strings.Join(strings.Fields(sql), " ")
	// Trim leading and trailing whitespace
	sql = strings.TrimSpace(sql)
	// Convert to lowercase for case-insensitive comparison
	return strings.ToLower(sql)
}
