package models

import (
	"errors"
	"math"
	"sort"
	"time"
)

type Scores struct {
	ID                     int64     `json:"id"`
	BrierScore             float64   `json:"brier_score"`
	Log2Score              float64   `json:"log2_score"`
	LogNScore              float64   `json:"logn_score"`
	BrierScoreTimeWeighted float64   `json:"brier_score_time_weighted"`
	Log2ScoreTimeWeighted  float64   `json:"log2_score_time_weighted"`
	LogNScoreTimeWeighted  float64   `json:"logn_score_time_weighted"`
	UserID                 int64     `json:"user_id"`
	ForecastID             int64     `json:"forecast_id"`
	CreatedAt              time.Time `json:"created"`
}

// Base struct for common score fields
type ScoreMetrics struct {
	BrierScore             float64 `json:"brier_score"`
	Log2Score              float64 `json:"log2_score"`
	LogNScore              float64 `json:"logn_score"`
	BrierScoreTimeWeighted float64 `json:"brier_score_time_weighted"`
	Log2ScoreTimeWeighted  float64 `json:"log2_score_time_weighted"`
	LogNScoreTimeWeighted  float64 `json:"logn_score_time_weighted"`
}

type ScoreFilters struct {
	UserID        *int64
	ForecastID    *int64
	Category      *string
	GroupByUserID *bool
	StartDate     *time.Time
	EndDate       *time.Time
}

// Overall platform averages
type OverallScores struct {
	ScoreMetrics
	TotalUsers     int `json:"total_users"`
	TotalForecasts int `json:"total_forecasts"`
}

// Per-category averages
type CategoryScores struct {
	ScoreMetrics
	Category       string `json:"category"`
	TotalUsers     int    `json:"total_users"`
	TotalForecasts int    `json:"total_forecasts"`
}

// Per-user averages
type UserScores struct {
	ScoreMetrics
	UserID         int64 `json:"user_id"`
	TotalForecasts int   `json:"total_forecasts"`
}

// Per-user-per-category averages
type UserCategoryScores struct {
	ScoreMetrics
	UserID         int64  `json:"user_id"`
	Category       string `json:"category"`
	TotalForecasts int    `json:"total_forecasts"`
}

type TimePoint struct {
	PointForecast float64
	CreatedAt     time.Time
}

func CalcForecastScore(points []TimePoint, outcome bool, userID int64, forecastID int64, forecastCreatedAt time.Time, forecastClosingDate *time.Time, forecastResolvedAt *time.Time) (Scores, error) {
	if len(points) == 0 {
		return Scores{}, errors.New("no probabilities provided")
	}

	// Sort points by CreatedAt to ensure correct time weighting
	sort.Slice(points, func(i, j int) bool {
		return points[i].CreatedAt.Before(points[j].CreatedAt)
	})

	var closeDate time.Time
	if forecastClosingDate != nil && forecastClosingDate.Before(*forecastResolvedAt) {
		closeDate = *forecastClosingDate
	} else {
		closeDate = *forecastResolvedAt
	}

	var brierSum, log2Sum, logNSum float64
	var brierSumTimeWeighted, log2SumTimeWeighted, logNSumTimeWeighted float64
	pointsCount := float64(len(points))
	firstPointCreatedAt := points[0].CreatedAt
	totalTimeInForecast := closeDate.Sub(firstPointCreatedAt).Seconds()

	// Edge case: if forecast was created and resolved at the same time (or very close),
	// fall back to naive (equal-weighted) scoring for time-weighted scores
	useTimeWeighting := totalTimeInForecast > 1.0 // At least 1 second open

	for i, point := range points {
		if point.PointForecast <= 0.0 || point.PointForecast >= 1.0 {
			return Scores{}, errors.New("point forecasts must be within 0 and 1")
		}

		// Calculate time weight
		var timeWeight float64
		if useTimeWeighting {
			// Calculate how long this prediction was held
			var duration float64
			if i < len(points)-1 {
				duration = points[i+1].CreatedAt.Sub(point.CreatedAt).Seconds()
			} else {
				duration = closeDate.Sub(point.CreatedAt).Seconds()
			}
			timeWeight = duration / totalTimeInForecast
		} else {
			// Fall back to equal weighting if no time elapsed
			timeWeight = 1.0 / pointsCount
		}

		if outcome {
			brierSum += math.Pow(point.PointForecast-1, 2)
			logNSum += math.Log(point.PointForecast)
			log2Sum += math.Log2(point.PointForecast)
			brierSumTimeWeighted += math.Pow(point.PointForecast-1, 2) * timeWeight
			logNSumTimeWeighted += math.Log(point.PointForecast) * timeWeight
			log2SumTimeWeighted += math.Log2(point.PointForecast) * timeWeight
		} else {
			brierSum += math.Pow(point.PointForecast, 2)
			logNSum += math.Log(1 - point.PointForecast)
			log2Sum += math.Log2(1 - point.PointForecast)
			brierSumTimeWeighted += math.Pow(point.PointForecast, 2) * timeWeight
			logNSumTimeWeighted += math.Log(1-point.PointForecast) * timeWeight
			log2SumTimeWeighted += math.Log2(1-point.PointForecast) * timeWeight
		}
	}

	return Scores{
		BrierScore:             brierSum / pointsCount,
		Log2Score:              log2Sum / pointsCount,
		LogNScore:              logNSum / pointsCount,
		BrierScoreTimeWeighted: brierSumTimeWeighted,
		LogNScoreTimeWeighted:  logNSumTimeWeighted,
		Log2ScoreTimeWeighted:  log2SumTimeWeighted,
		UserID:                 userID,
		ForecastID:             forecastID,
		CreatedAt:              time.Now(),
	}, nil
}
