package models

import (
	"encoding/json"
	"reflect"
)

// REQUEST TYPES
type PerplexityRequest struct {
	Messages  []Message `json:"messages" validate:"required,dive"`
	Model     string    `json:"model" validate:"required"`
	MaxTokens int       `json:"max_tokens" validate:"gt=0"`
}

// RESPONSE TYPES
// type for typical messages
type Message struct {
	Role    string `json:"role" validate:"required,oneof=system user assistant"`
	Content string `json:"content"`
}

// type for the token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// choices
type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
	Delta        Message `json:"delta"`
}

// full response
type PerplexityResponse struct {
	ID               string          `json:"id"`
	Model            string          `json:"model"`
	Created          int             `json:"created"`
	Usage            Usage           `json:"usage"`
	Object           string          `json:"object"`
	Choices          []Choice        `json:"choices"`
	SearchResults    *[]SearchResult `json:"search_results,omitempty"`
	RelatedQuestions *[]string       `json:"related_questions,omitempty"`
}

// single search result
type SearchResult struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Date        *string `json:"date,omitempty"`
	LastUpdated *string `json:"last_updated,omitempty"`
}

// string representation of the SearchResult
func (sr *SearchResult) String() string {
	if sr == nil {
		return ""
	}

	result := sr.Title
	if sr.URL != "" {
		result += " (" + sr.URL + ")"
	}
	if sr.Date != nil {
		result += " - " + *sr.Date
	}
	if sr.LastUpdated != nil {
		result += " (updated: " + *sr.LastUpdated + ")"
	}

	return result
}

// String returns a string representation of the PerplexityResponse.
func (r *PerplexityResponse) String() string {
	if r == nil {
		return ""
	}
	if reflect.DeepEqual(r, &PerplexityResponse{}) {
		return ""
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// GetLastContent returns the last content of the PerplexityResponse.
func (r *PerplexityResponse) GetLastContent() string {
	if len(r.Choices) == 0 {
		return ""
	}
	return r.Choices[len(r.Choices)-1].Message.Content
}

// GetSearchResults returns the search results of the PerplexityResponse.
func (r *PerplexityResponse) GetSearchResults() []SearchResult {
	if r.SearchResults == nil {
		return []SearchResult{}
	}
	return *r.SearchResults
}
