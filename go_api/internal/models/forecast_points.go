package models

import (
	"time"
)

type ForecastPoint struct {
	ID            int64     `json:"id"`
	ForecastID    int64     `json:"forecast_id"`
	PointForecast float64   `json:"point_forecast"`
	UpperCI       float64   `json:"upper_ci"`
	LowerCI       float64   `json:"lower_ci"`
	CreatedAt     time.Time `json:"created"`
	Reason        string    `json:"reason"`
}
