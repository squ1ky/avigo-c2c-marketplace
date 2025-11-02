package kafka

import "time"

type EventType string

const (
	EventListingCreated   EventType = "listing.created"
	EventListingUpdated   EventType = "listing.updated"
	EventListingDeleted   EventType = "listing.deleted"
	EventListingPublished EventType = "listing.published"
)

type ListingEvent struct {
	EventID   string                 `json:"event_id"`
	EventType EventType              `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id"`
	ListingID string                 `json:"listing_id"`
	Data      map[string]interface{} `json:"data"`
}
