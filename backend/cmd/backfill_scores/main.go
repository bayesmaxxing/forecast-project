package main

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// This script recalculates all scores for existing score records
// Used to fix time-weighting bug where denominator now uses user's first point time
// Run with: go run cmd/backfill_scores/main.go

func main() {
	log.Println("Starting score backfill...")

	// Initialize database connection
	db, err := database.NewDB(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Get all scores that need backfilling (where time_weighted columns are NULL)
	query := `
		SELECT
			s.id,
			s.user_id,
			s.forecast_id,
			f.created as forecast_created,
			f.resolved as forecast_resolved,
			f.closing_date as forecast_closing
		FROM scores s
		JOIN forecasts f ON s.forecast_id = f.id
		WHERE f.resolved IS NOT NULL
		ORDER BY s.id
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to query scores: %v", err)
	}
	defer rows.Close()

	type ScoreToBackfill struct {
		ScoreID          int64
		UserID           int64
		ForecastID       int64
		ForecastCreated  time.Time
		ForecastResolved time.Time
		ForecastClosing  *time.Time
	}

	var scoresToBackfill []ScoreToBackfill
	for rows.Next() {
		var s ScoreToBackfill
		if err := rows.Scan(&s.ScoreID, &s.UserID, &s.ForecastID, &s.ForecastCreated, &s.ForecastResolved, &s.ForecastClosing); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		scoresToBackfill = append(scoresToBackfill, s)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	log.Printf("Found %d scores to backfill", len(scoresToBackfill))

	if len(scoresToBackfill) == 0 {
		log.Println("No scores to backfill. Exiting.")
		return
	}

	// Ask for confirmation
	fmt.Printf("About to backfill %d scores. Continue? (y/n): ", len(scoresToBackfill))
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		log.Println("Backfill cancelled.")
		return
	}

	successCount := 0
	errorCount := 0

	// Process each score
	for i, scoreToBackfill := range scoresToBackfill {
		if i%10 == 0 {
			log.Printf("Progress: %d/%d", i, len(scoresToBackfill))
		}

		// Get all points for this forecast and user
		pointsQuery := `
			SELECT point_forecast, created
			FROM points
			WHERE forecast_id = $1 AND user_id = $2
			ORDER BY created ASC
		`

		pointRows, err := db.QueryContext(ctx, pointsQuery, scoreToBackfill.ForecastID, scoreToBackfill.UserID)
		if err != nil {
			log.Printf("Error fetching points for score %d: %v", scoreToBackfill.ScoreID, err)
			errorCount++
			continue
		}

		var points []models.TimePoint
		for pointRows.Next() {
			var p models.TimePoint
			if err := pointRows.Scan(&p.PointForecast, &p.CreatedAt); err != nil {
				log.Printf("Error scanning point for score %d: %v", scoreToBackfill.ScoreID, err)
				pointRows.Close()
				errorCount++
				continue
			}
			points = append(points, p)
		}
		pointRows.Close()

		if len(points) == 0 {
			log.Printf("No points found for score %d, skipping", scoreToBackfill.ScoreID)
			errorCount++
			continue
		}

		// Get the resolution outcome
		var resolutionStr *string
		resolutionQuery := `SELECT resolution FROM forecasts WHERE id = $1`
		err = db.QueryRowContext(ctx, resolutionQuery, scoreToBackfill.ForecastID).Scan(&resolutionStr)
		if err != nil {
			log.Printf("Error fetching resolution for forecast %d: %v", scoreToBackfill.ForecastID, err)
			errorCount++
			continue
		}

		if resolutionStr == nil || *resolutionStr == "-" {
			log.Printf("Forecast %d has no valid resolution, skipping score %d", scoreToBackfill.ForecastID, scoreToBackfill.ScoreID)
			errorCount++
			continue
		}

		outcome := *resolutionStr == "1"

		// Recalculate scores
		recalculatedScore, err := models.CalcForecastScore(
			points,
			outcome,
			scoreToBackfill.UserID,
			scoreToBackfill.ForecastID,
			scoreToBackfill.ForecastCreated,
			scoreToBackfill.ForecastClosing,
			&scoreToBackfill.ForecastResolved,
		)
		if err != nil {
			log.Printf("Error calculating score for score %d: %v", scoreToBackfill.ScoreID, err)
			errorCount++
			continue
		}

		// Update all score columns
		updateQuery := `
			UPDATE scores
			SET
				brier_score = $1,
				log2_score = $2,
				logn_score = $3,
				brier_score_time_weighted = $4,
				log2_score_time_weighted = $5,
				logn_score_time_weighted = $6
			WHERE id = $7
		`

		_, err = db.ExecContext(ctx, updateQuery,
			recalculatedScore.BrierScore,
			recalculatedScore.Log2Score,
			recalculatedScore.LogNScore,
			recalculatedScore.BrierScoreTimeWeighted,
			recalculatedScore.Log2ScoreTimeWeighted,
			recalculatedScore.LogNScoreTimeWeighted,
			scoreToBackfill.ScoreID,
		)
		if err != nil {
			log.Printf("Error updating score %d: %v", scoreToBackfill.ScoreID, err)
			errorCount++
			continue
		}

		successCount++
	}

	log.Printf("Backfill complete!")
	log.Printf("Success: %d", successCount)
	log.Printf("Errors: %d", errorCount)

	if errorCount > 0 {
		os.Exit(1)
	}
}
