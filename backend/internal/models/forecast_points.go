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
	UserName      *string   `json:"user_name,omitempty"`
}

type PointFilters struct {
	UserID             *int64
	ForecastID         *int64
	Date               *time.Time
	DistinctOnForecast *bool
	OrderByForecastID  *bool
	CreatedDirection   *string
}
