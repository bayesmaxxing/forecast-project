package models

import (
	"errors"
	"math"
	"time"
)

type Scores struct {
	ID         int64     `json:"id"`
	BrierScore float64   `json:"brier_score"`
	Log2Score  float64   `json:"log2_score"`
	LogNScore  float64   `json:"logN_score"`
	UserID     int64     `json:"user_id"`
	ForecastID int64     `json:"forecast_id"`
	CreatedAt  time.Time `json:"created"`
}

// Base struct for common score fields
type ScoreMetrics struct {
	BrierScore float64 `json:"brier_score"`
	Log2Score  float64 `json:"log2_score"`
	LogNScore  float64 `json:"logn_score"`
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

func CalcForecastScore(probabilities []float64, outcome bool, userID int64, forecastID int64) (Scores, error) {
	if len(probabilities) == 0 {
		return Scores{}, errors.New("no probabilities provided")
	}

	var brierSum, log2Sum, logNSum float64
	points := float64(len(probabilities))

	for _, prob := range probabilities {
		if prob <= 0.0 || prob >= 1.0 {
			return Scores{}, errors.New("probs must be within 0 and 1")
		}

		if outcome {
			brierSum += math.Pow(prob-1, 2)
			logNSum += math.Log(prob)
			log2Sum += math.Log2(prob)
		} else {
			brierSum += math.Pow(prob, 2)
			logNSum += math.Log(1 - prob)
			log2Sum += math.Log2(1 - prob)
		}
	}

	return Scores{
		BrierScore: brierSum / points,
		Log2Score:  log2Sum / points,
		LogNScore:  logNSum / points,
		UserID:     userID,
		ForecastID: forecastID,
		CreatedAt:  time.Now(),
	}, nil
}
