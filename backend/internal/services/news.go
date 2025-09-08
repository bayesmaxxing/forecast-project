package services

import (
	"backend/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type NewsService struct {
}

func NewNewsService() *NewsService {
	return &NewsService{}
}

func (s *NewsService) GetNews(ctx context.Context, query string) (*models.PerplexityResponse, error) {

	r := &models.PerplexityResponse{}
	request := models.PerplexityRequest{
		Messages: []models.Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		Model:     "sonar-pro",
		MaxTokens: 2000,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.perplexity.ai/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("PERPLEXITY_API_KEY"))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check return status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("unauthorized")
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code (%d) and cannot read response: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("unexpected status code (%d): %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w - body response=%s", err, string(body))
	}
	return r, nil
}
