package models

import "time"

// CalibrationBucket represents a single probability bucket with calibration data
type CalibrationBucket struct {
	BucketStart     float64 `json:"bucket_start"`
	BucketEnd       float64 `json:"bucket_end"`
	PredictionCount int     `json:"prediction_count"`
	AvgPrediction   float64 `json:"avg_prediction"`
	ActualRate      float64 `json:"actual_rate"`
}

// CalibrationData contains overall calibration data
type CalibrationData struct {
	Buckets          []CalibrationBucket `json:"buckets"`
	TotalPredictions int                 `json:"total_predictions"`
	TotalForecasts   int                 `json:"total_forecasts"`
}

// UserCalibrationData contains per-user calibration data
type UserCalibrationData struct {
	UserID           int64               `json:"user_id"`
	Buckets          []CalibrationBucket `json:"buckets"`
	TotalPredictions int                 `json:"total_predictions"`
	TotalForecasts   int                 `json:"total_forecasts"`
}

// CalibrationFilters contains filter options for calibration queries
type CalibrationFilters struct {
	UserID    *int64
	Category  *string
	StartDate *time.Time
	EndDate   *time.Time
}
