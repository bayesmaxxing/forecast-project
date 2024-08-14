package models

import (
	"time"
)

type Forecast struct {
	ID                 int64      `json:"id"`
	Question           string     `json:"question"`
	Category           string     `json:"category"`
	CreatedAt          time.Time  `json:"created"`
	ResolutionCriteria string     `json:"resolution_criteria"`
	Resolution         *string    `json:"resolution"`
	ResolvedAt         *time.Time `json:"resolved"`
	BrierScore         *float64   `json:"brier_score"`
	Log2Score          *float64   `json:"log2_score"`
	LogNScore          *float64   `json:"logn_score"`
	ResolutionComment  *string    `json:"comment"`
}

// Check if forecast has resolved
func (f *Forecast) IsResolved() bool {
	return f.ResolvedAt != nil
}

func (f *Forecast) Resolve(resolution string, comment string) {
	now := time.Now()
	f.ResolvedAt = &now
	f.Resolution = &resolution
	f.ResolutionComment = &comment

	//import brier_score function and log2score and lognscore here
	// also need to get the forecast_points in here...
}
