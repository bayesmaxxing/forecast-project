package models

import (
	"errors"
	"go_api/internal/utils"
	"time"
)

type Forecast struct {
	ID                 int64      `json:"id"`
	Question           string     `json:"question"`
	Category           string     `json:"category"`
	CreatedAt          time.Time  `json:"created"`
	ResolutionCriteria string     `json:"resolution_criteria"`
	Resolution         *string    `json:"resolution,omitempty"`
	ResolvedAt         *time.Time `json:"resolved,omitempty"`
	BrierScore         *float64   `json:"brier_score,omitempty"`
	Log2Score          *float64   `json:"log2_score,omitempty"`
	LogNScore          *float64   `json:"logn_score,omitempty"`
	ResolutionComment  *string    `json:"comment,omitempty"`
}

// Check if forecast has resolved
func (f *Forecast) IsResolved() bool {
	return f.ResolvedAt != nil
}

func (f *Forecast) Resolve(resolution string, comment string, probabilities []float64) error {
	if f.ResolvedAt != nil {
		return errors.New("forecast has already been resolved")
	}

	if len(probabilities) == 0 {
		return errors.New("no probabilities supplies")
	}

	now := time.Now()
	f.ResolvedAt = &now
	f.Resolution = &resolution
	f.ResolutionComment = &comment

	outcome := resolution == "1"

	scores, err := utils.CalcForecastScores(probabilities, outcome)

	if err != nil {
		return err
	}

	f.BrierScore = &scores.BrierScore
	f.Log2Score = &scores.Log2Score
	f.LogNScore = &scores.LogNScore
	return nil
}
