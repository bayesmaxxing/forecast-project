package models

import (
	"time"
)

type Blogpost struct {
	ID               int64     `json:"post_id"`
	Title            string    `json:"title"`
	Post             string    `json:"post"`
	CreatedAt        time.Time `json:"created"`
	Summary          string    `json:"summary"`
	Slug             string    `json:"slug"`
	RelatedForecasts []int64   `json:"related_forecasts"`
}
