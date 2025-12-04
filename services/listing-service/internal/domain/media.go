package domain

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type ListingMedia struct {
	ID        uuid.UUID `db:"id"`
	ListingID uuid.UUID `db:"listing_id"`
	FileURL   string    `db:"file_url"` // MinIO path
	FileType  MediaType `db:"file_type"`
	MimeType  string    `db:"mime_type"` // image/jpeg, image/png
	Order     int       `db:"order"`
	CreatedAt time.Time `db:"created_at"`
}

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

func DetectMediaType(mimeType string) MediaType {
	if strings.HasPrefix(mimeType, "video") {
		return MediaTypeVideo
	}
	return MediaTypeImage
}
