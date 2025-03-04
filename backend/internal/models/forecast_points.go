package models

import (
	"time"
)

type ForecastPoint struct {
	ID            int64     `json:"id"`
	ForecastID    int64     `json:"forecast_id"`
	PointForecast float64   `json:"point_forecast"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created"`
	UserID        int64     `json:"user_id"`
}
