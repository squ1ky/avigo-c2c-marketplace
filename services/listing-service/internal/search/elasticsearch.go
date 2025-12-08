package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
)

type Client struct {
	es    *elasticsearch.Client
	index string
}

type ListingDocument struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	IsSold      bool      `json:"is_sold"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   string    `json:"created_at"`
}

type ListingSearchParams struct {
	Query      string
	From       int
	Size       int
	CategoryID *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
}

type ListingSearchResult struct {
	IDs   []uuid.UUID
	Total int64
}

func NewClient(cfg config.ElasticsearchConfig) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init elasticsearch client: %w", err)
	}

	return &Client{es: es, index: cfg.Index}, nil
}

func (c *Client) IndexListing(ctx context.Context, doc ListingDocument) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal listing doc: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: doc.ID.String(),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("elasticsearch index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("elasticsearch index failed: status %s", res.Status())
	}

	return nil
}

func (c *Client) DeleteListing(ctx context.Context, id uuid.UUID) error {
	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: id.String(),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("elasticsearch delete request failed: %w", err)
	}
	defer res.Body.Close()

	// 404 is not fatal for idempotency
	if res.StatusCode >= http.StatusBadRequest && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("elasticsearch delete failed: status %s", res.Status())
	}

	return nil
}

func (c *Client) SearchListings(ctx context.Context, params ListingSearchParams) (*ListingSearchResult, error) {
	boolQuery := map[string]interface{}{
		"filter": []interface{}{
			map[string]interface{}{"term": map[string]interface{}{"status": "active"}},
			map[string]interface{}{"term": map[string]interface{}{"is_sold": false}},
		},
	}

	if params.Query != "" {
		boolQuery["must"] = []interface{}{
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  params.Query,
					"fields": []string{"title^3", "description", "tags^2"},
				},
			},
		}
	} else {
		boolQuery["must"] = []interface{}{
			map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}

	if params.CategoryID != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}),
			map[string]interface{}{"term": map[string]interface{}{"category_id": params.CategoryID.String()}},
		)
	}

	if params.MinPrice != nil || params.MaxPrice != nil {
		priceRange := map[string]interface{}{}
		if params.MinPrice != nil {
			priceRange["gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			priceRange["lte"] = *params.MaxPrice
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}),
			map[string]interface{}{"range": map[string]interface{}{"price": priceRange}},
		)
	}

	searchBody := map[string]interface{}{
		"from": params.From,
		"size": params.Size,
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"sort": []interface{}{
			map[string]interface{}{"_score": map[string]interface{}{"order": "desc"}},
			map[string]interface{}{"created_at": map[string]interface{}{"order": "desc"}},
		},
	}

	body, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search body: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(c.index),
		c.es.Search.WithBody(bytes.NewReader(body)),
		c.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("elasticsearch search failed: status %s", res.Status())
	}

	var parsed struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	result := &ListingSearchResult{
		IDs:   make([]uuid.UUID, 0, len(parsed.Hits.Hits)),
		Total: parsed.Hits.Total.Value,
	}

	for _, hit := range parsed.Hits.Hits {
		id, err := uuid.Parse(hit.ID)
		if err != nil {
			slog.Warn("invalid listing id in elasticsearch result", "id", hit.ID, "error", err)
			continue
		}
		result.IDs = append(result.IDs, id)
	}

	return result, nil
}
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
)

type Client struct {
	es    *elasticsearch.Client
	index string
}

type ListingDocument struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	IsSold      bool      `json:"is_sold"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   string    `json:"created_at"`
}

type ListingSearchParams struct {
	Query      string
	From       int
	Size       int
	CategoryID *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
}

type ListingSearchResult struct {
	IDs   []uuid.UUID
	Total int64
}

func NewClient(cfg config.ElasticsearchConfig) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init elasticsearch client: %w", err)
	}

	return &Client{es: es, index: cfg.Index}, nil
}

func (c *Client) IndexListing(ctx context.Context, doc ListingDocument) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal listing doc: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: doc.ID.String(),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("elasticsearch index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("elasticsearch index failed: status %s", res.Status())
	}

	return nil
}

func (c *Client) DeleteListing(ctx context.Context, id uuid.UUID) error {
	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: id.String(),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("elasticsearch delete request failed: %w", err)
	}
	defer res.Body.Close()

	// 404 is not fatal for idempotency
	if res.StatusCode >= http.StatusBadRequest && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("elasticsearch delete failed: status %s", res.Status())
	}

	return nil
}

func (c *Client) SearchListings(ctx context.Context, params ListingSearchParams) (*ListingSearchResult, error) {
	boolQuery := map[string]interface{}{
		"filter": []interface{}{
			map[string]interface{}{"term": map[string]interface{}{"status": "active"}},
			map[string]interface{}{"term": map[string]interface{}{"is_sold": false}},
		},
	}

	if params.Query != "" {
		boolQuery["must"] = []interface{}{
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  params.Query,
					"fields": []string{"title^3", "description", "tags^2"},
				},
			},
		}
	} else {
		boolQuery["must"] = []interface{}{
			map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}

	if params.CategoryID != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}),
			map[string]interface{}{"term": map[string]interface{}{"category_id": params.CategoryID.String()}},
		)
	}

	if params.MinPrice != nil || params.MaxPrice != nil {
		priceRange := map[string]interface{}{}
		if params.MinPrice != nil {
			priceRange["gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			priceRange["lte"] = *params.MaxPrice
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}),
			map[string]interface{}{"range": map[string]interface{}{"price": priceRange}},
		)
	}

	searchBody := map[string]interface{}{
		"from": params.From,
		"size": params.Size,
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"sort": []interface{}{
			map[string]interface{}{"_score": map[string]interface{}{"order": "desc"}},
			map[string]interface{}{"created_at": map[string]interface{}{"order": "desc"}},
		},
	}

	body, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search body: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(c.index),
		c.es.Search.WithBody(bytes.NewReader(body)),
		c.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("elasticsearch search failed: status %s", res.Status())
	}

	var parsed struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	result := &ListingSearchResult{
		IDs:   make([]uuid.UUID, 0, len(parsed.Hits.Hits)),
		Total: parsed.Hits.Total.Value,
	}

	for _, hit := range parsed.Hits.Hits {
		id, err := uuid.Parse(hit.ID)
		if err != nil {
			slog.Warn("invalid listing id in elasticsearch result", "id", hit.ID, "error", err)
			continue
		}
		result.IDs = append(result.IDs, id)
	}

	return result, nil
}

