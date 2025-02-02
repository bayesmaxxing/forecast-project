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
