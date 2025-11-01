package models

import (
	"time"
)

type Forecast struct {
	ID                 int64      `json:"id"`
	Question           string     `json:"question"`
	Category           string     `json:"category"`
	CreatedAt          time.Time  `json:"created"`
	UserID             int64      `json:"user_id"`
	ResolutionCriteria string     `json:"resolution_criteria"`
	ClosingDate        *time.Time `json:"closing_date,omitempty"`
	Resolution         *string    `json:"resolution,omitempty"`
	ResolvedAt         *time.Time `json:"resolved,omitempty"`
	ResolutionComment  *string    `json:"comment,omitempty"`
}

// Check if forecast has resolved
func (f *Forecast) IsResolved() bool {
	return f.ResolvedAt != nil
}
