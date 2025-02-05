package models

import (
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
