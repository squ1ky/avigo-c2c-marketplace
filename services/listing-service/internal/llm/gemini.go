package llm

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
)

var ErrNotConfigured = errors.New("gemini api key is not configured")

type TagSuggestionRequest struct {
	Title        string
	CategoryName string
}

type GeminiClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewGeminiClient(cfg config.GeminiConfig) *GeminiClient {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}

	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = "gemini-2.5-flash"
	}

	return &GeminiClient{
		apiKey:  strings.TrimSpace(cfg.APIKey),
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *GeminiClient) SuggestTags(ctx context.Context, input TagSuggestionRequest) ([]string, error) {
	if c.apiKey == "" {
		return nil, ErrNotConfigured
	}

	reqBody := geminiGenerateRequest{
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: buildTagPrompt(input)},
				},
			},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:      0.2,
			MaxOutputTokens:  256,
			ThinkingConfig:   geminiThinkingConfig{ThinkingBudget: 0},
			ResponseMimeType: "application/json",
			ResponseSchema: geminiSchema{
				Type: "OBJECT",
				Properties: map[string]geminiSchema{
					"tags": {
						Type: "ARRAY",
						Items: &geminiSchema{
							Type: "STRING",
						},
					},
				},
				Required: []string{"tags"},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent", c.baseURL, c.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to build gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("failed to read gemini response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var genResp geminiGenerateResponse
	if err := json.Unmarshal(respBody, &genResp); err != nil {
		return nil, fmt.Errorf("failed to decode gemini response: %w", err)
	}

	text := genResp.text()
	if text == "" {
		return nil, errors.New("gemini returned empty tag suggestion")
	}

	var parsed struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(cleanJSONText(text)), &parsed); err != nil {
		return nil, fmt.Errorf("failed to decode gemini tags json: %w", err)
	}

	return parsed.Tags, nil
}

func buildTagPrompt(input TagSuggestionRequest) string {
	var b strings.Builder
	b.WriteString("You generate listing tags for a C2C marketplace.\n")
	b.WriteString("Return 5-10 short tags in Russian.\n")
	b.WriteString("Use common Latin brand/model variants only when they are clearly present in the title.\n")
	b.WriteString("Do not invent specs, price, city, item condition, or promotional words.\n")
	b.WriteString("Return only valid JSON, without markdown or explanation.\n")
	b.WriteString("The exact JSON shape is: {\"tags\":[\"tag\"]}.\n\n")
	b.WriteString("Title: ")
	b.WriteString(input.Title)
	if input.CategoryName != "" {
		b.WriteString("\nCategory: ")
		b.WriteString(input.CategoryName)
	}
	return b.String()
}

func cleanJSONText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end >= start {
		return strings.TrimSpace(text[start : end+1])
	}

	return text
}

type geminiGenerateRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature      float64              `json:"temperature,omitempty"`
	MaxOutputTokens  int                  `json:"maxOutputTokens,omitempty"`
	ThinkingConfig   geminiThinkingConfig `json:"thinkingConfig,omitempty"`
	ResponseMimeType string               `json:"responseMimeType,omitempty"`
	ResponseSchema   geminiSchema         `json:"responseSchema,omitempty"`
}

type geminiThinkingConfig struct {
	ThinkingBudget int `json:"thinkingBudget"`
}

type geminiSchema struct {
	Type       string                  `json:"type,omitempty"`
	Properties map[string]geminiSchema `json:"properties,omitempty"`
	Items      *geminiSchema           `json:"items,omitempty"`
	Required   []string                `json:"required,omitempty"`
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

func (r geminiGenerateResponse) text() string {
	var b strings.Builder
	for _, candidate := range r.Candidates {
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.Text) != "" {
				if b.Len() > 0 {
					b.WriteString("\n")
				}
				b.WriteString(part.Text)
			}
		}
	}
	return strings.TrimSpace(b.String())
}
