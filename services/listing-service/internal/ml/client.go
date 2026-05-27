package ml

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/llm"
)

var ErrNotConfigured = errors.New("ml service url is not configured")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg config.MLServiceConfig) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 3 * time.Second
	}

	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(cfg.URL), "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) SuggestTags(ctx context.Context, input llm.TagSuggestionRequest) ([]string, error) {
	if c.baseURL == "" {
		return nil, ErrNotConfigured
	}

	reqBody := predictTagsRequest{
		Title:    input.Title,
		Category: input.CategoryName,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ml request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/predict-tags", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to build ml request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ml request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("failed to read ml response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ml service returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var parsed predictTagsResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("failed to decode ml response: %w", err)
	}

	return parsed.Tags, nil
}

type predictTagsRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}

type predictTagsResponse struct {
	Tags []string `json:"tags"`
}
