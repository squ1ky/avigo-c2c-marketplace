package domain

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	ParentID  uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Level     int       `json:"level" db:"level"` // 0 - root, 1 - subcategory, ...
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AttributeType string

const (
	AttributeTypeString  AttributeType = "string"
	AttributeTypeNumber  AttributeType = "number"
	AttributeTypeBoolean AttributeType = "boolean"
	AttributeTypeSelect  AttributeType = "select"
	AttributeTypeMulti   AttributeType = "multi"
)

type CategoryAttribute struct {
	ID         uuid.UUID     `json:"id" db:"id"`
	CategoryID uuid.UUID     `json:"category_id" db:"category_id"`
	Name       string        `json:"name" db:"name"`           // "Диагональ экрана"
	Key        string        `json:"key" db:"key"`             // "screen_diagonal"
	Type       AttributeType `json:"type" db:"type"`           // string, boolean, number
	Unit       *string       `json:"unit,omitempty" db:"unit"` // "дюйм", "см", "кг"
	Order      int           `json:"order" db:"order"`
	Options    []string      `json:"options,omitempty"` // type=select: ["Wood", "Stone"]
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
}
