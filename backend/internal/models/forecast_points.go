package models

import (
	"time"
)

type ForecastPoint struct {
	ID            int64     `json:"update_id"`
	ForecastID    int64     `json:"forecast_id"`
	PointForecast float64   `json:"point_forecast"`
	UpperCI       float64   `json:"upper_ci"`
	LowerCI       float64   `json:"lower_ci"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created"`
	UserID        int64     `json:"user_id"`
}
