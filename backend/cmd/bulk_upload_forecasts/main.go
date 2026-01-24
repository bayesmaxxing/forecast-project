package main

import (
	"backend/internal/database"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// This script bulk uploads forecasts from a CSV file
// CSV format: question,category,resolution_criteria,user_id,closing_date
// - question: The forecast question (required)
// - category: The forecast category (required)
// - resolution_criteria: Criteria for resolution (required)
// - user_id: The user ID who owns the forecast (required)
// - closing_date: When the forecast closes, RFC3339 format e.g. 2025-12-31T23:59:59Z (optional)
//
// Run with: go run cmd/bulk_upload_forecasts/main.go <csv_file_path>

type ForecastCSV struct {
	Question           string
	Category           string
	ResolutionCriteria string
	UserID             int64
	ClosingDate        *time.Time
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/bulk_upload_forecasts/main.go <csv_file_path>")
	}

	csvPath := os.Args[1]
	log.Printf("Reading forecasts from: %s", csvPath)

	// Read and parse CSV
	forecasts, err := parseCSV(csvPath)
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	log.Printf("Found %d forecasts to upload", len(forecasts))

	if len(forecasts) == 0 {
		log.Println("No forecasts to upload. Exiting.")
		return
	}

	// Initialize database connection
	db, err := database.NewDB(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ask for confirmation
	fmt.Printf("About to upload %d forecasts. Continue? (y/n): ", len(forecasts))
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		log.Println("Upload cancelled.")
		return
	}

	ctx := context.Background()
	successCount := 0
	errorCount := 0

	// Insert forecasts
	for i, f := range forecasts {
		if i%10 == 0 && i > 0 {
			log.Printf("Progress: %d/%d", i, len(forecasts))
		}

		id, err := insertForecast(ctx, db, f)
		if err != nil {
			log.Printf("Error inserting forecast %d (%s): %v", i+1, truncate(f.Question, 50), err)
			errorCount++
			continue
		}

		log.Printf("Created forecast ID %d: %s", id, truncate(f.Question, 50))
		successCount++
	}

	log.Println("Upload complete!")
	log.Printf("Success: %d", successCount)
	log.Printf("Errors: %d", errorCount)

	if errorCount > 0 {
		os.Exit(1)
	}
}

func parseCSV(path string) ([]ForecastCSV, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Map header columns to indices
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	// Validate required columns
	requiredCols := []string{"question", "category", "resolution_criteria", "user_id"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var forecasts []ForecastCSV
	lineNum := 1 // Header is line 1

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNum, err)
		}

		// Parse required fields
		question := record[colIndex["question"]]
		category := record[colIndex["category"]]
		resolutionCriteria := record[colIndex["resolution_criteria"]]
		userIDStr := record[colIndex["user_id"]]

		if question == "" {
			return nil, fmt.Errorf("line %d: question is required", lineNum)
		}
		if category == "" {
			return nil, fmt.Errorf("line %d: category is required", lineNum)
		}
		if resolutionCriteria == "" {
			return nil, fmt.Errorf("line %d: resolution_criteria is required", lineNum)
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid user_id '%s': %w", lineNum, userIDStr, err)
		}

		// Parse optional closing_date
		var closingDate *time.Time
		if idx, ok := colIndex["closing_date"]; ok && idx < len(record) && record[idx] != "" {
			t, err := time.Parse(time.RFC3339, record[idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid closing_date '%s' (expected RFC3339 format): %w", lineNum, record[idx], err)
			}
			closingDate = &t
		}

		forecasts = append(forecasts, ForecastCSV{
			Question:           question,
			Category:           category,
			ResolutionCriteria: resolutionCriteria,
			UserID:             userID,
			ClosingDate:        closingDate,
		})
	}

	return forecasts, nil
}

func insertForecast(ctx context.Context, db *database.DB, f ForecastCSV) (int64, error) {
	query := `
		INSERT INTO forecasts (question, category, created, user_id, resolution_criteria, closing_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := db.QueryRowContext(ctx, query,
		f.Question,
		f.Category,
		time.Now(),
		f.UserID,
		f.ResolutionCriteria,
		f.ClosingDate,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
