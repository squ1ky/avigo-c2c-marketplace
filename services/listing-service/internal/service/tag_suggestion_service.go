package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/llm"
	pgrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
)

var (
	ErrTagSuggestionsNotConfigured = errors.New("tag suggestions are not configured")
	ErrTagSuggestionsUnavailable   = errors.New("tag suggestions are unavailable")
)

type TagSuggestionProvider interface {
	SuggestTags(ctx context.Context, input llm.TagSuggestionRequest) ([]string, error)
}

type FallbackTagSuggestionProvider struct {
	primary  TagSuggestionProvider
	fallback TagSuggestionProvider
}

func NewFallbackTagSuggestionProvider(primary, fallback TagSuggestionProvider) *FallbackTagSuggestionProvider {
	return &FallbackTagSuggestionProvider{
		primary:  primary,
		fallback: fallback,
	}
}

func (p *FallbackTagSuggestionProvider) SuggestTags(ctx context.Context, input llm.TagSuggestionRequest) ([]string, error) {
	if p.primary != nil {
		tags, err := p.primary.SuggestTags(ctx, input)
		if err == nil && len(tags) > 0 {
			return tags, nil
		}
	}

	if p.fallback != nil {
		return p.fallback.SuggestTags(ctx, input)
	}

	return nil, ErrTagSuggestionsNotConfigured
}

type TagSuggestionService struct {
	provider    TagSuggestionProvider
	listingRepo *pgrepo.ListingRepository
}

func NewTagSuggestionService(provider TagSuggestionProvider, listingRepo *pgrepo.ListingRepository) *TagSuggestionService {
	return &TagSuggestionService{
		provider:    provider,
		listingRepo: listingRepo,
	}
}

func (s *TagSuggestionService) Suggest(ctx context.Context, input dto.SuggestTagsInput) (*dto.SuggestTagsResponse, error) {
	if s.provider == nil {
		return nil, ErrTagSuggestionsNotConfigured
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return &dto.SuggestTagsResponse{Tags: []string{}}, nil
	}

	categoryName := ""
	if input.CategoryID != "" {
		categoryID, err := uuid.Parse(input.CategoryID)
		if err == nil {
			if category, err := s.listingRepo.GetCategoryByID(ctx, categoryID); err == nil && category != nil {
				categoryName = category.Slug
			}
		}
	}

	tags, err := s.provider.SuggestTags(ctx, llm.TagSuggestionRequest{
		Title:        title,
		CategoryName: categoryName,
	})
	if err != nil {
		if errors.Is(err, llm.ErrNotConfigured) {
			return nil, ErrTagSuggestionsNotConfigured
		}
		return nil, fmt.Errorf("%w: %v", ErrTagSuggestionsUnavailable, err)
	}

	return &dto.SuggestTagsResponse{Tags: normalizeTags(tags, 10)}, nil
}

func normalizeTags(tags []string, limit int) []string {
	result := make([]string, 0, len(tags))
	seen := make(map[string]struct{}, len(tags))
	blocked := map[string]struct{}{
		"новый":              {},
		"новое":              {},
		"новая":              {},
		"б/у":                {},
		"бу":                 {},
		"отличное состояние": {},
	}

	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		tag = strings.Join(strings.Fields(tag), " ")
		if tag == "" || utf8.RuneCountInString(tag) > 32 {
			continue
		}
		if _, blocked := blocked[tag]; blocked {
			continue
		}
		if _, exists := seen[tag]; exists {
			continue
		}

		seen[tag] = struct{}{}
		result = append(result, tag)
		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result
}
