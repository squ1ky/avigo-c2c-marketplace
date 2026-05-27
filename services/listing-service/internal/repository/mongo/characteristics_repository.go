package mngrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CharacteristicsRepository struct {
	collection *mongo.Collection
}

func NewCharacteristicsRepository(db *mongo.Database) *CharacteristicsRepository {
	return &CharacteristicsRepository{
		collection: db.Collection("listing_characteristics"),
	}
}

func (r *CharacteristicsRepository) Create(ctx context.Context, chars *domain.ListingCharacteristics) error {
	chars.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, chars)
	if err != nil {
		return fmt.Errorf("failed to create characteristics: %w", err)
	}

	return nil
}

func (r *CharacteristicsRepository) Upsert(ctx context.Context, chars *domain.ListingCharacteristics) error {
	chars.UpdatedAt = time.Now()

	filter := bson.M{"listing_id": chars.ListingID}
	update := bson.M{
		"$set": bson.M{
			"characteristics": chars.Characteristics,
			"tags":            chars.Tags,
			"updated_at":      chars.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to upsert characteristics: %w", err)
	}

	return nil
}

func (r *CharacteristicsRepository) GetByListingID(ctx context.Context, listingID uuid.UUID) (*domain.ListingCharacteristics, error) {
	filter := bson.M{"listing_id": listingID}

	var chars domain.ListingCharacteristics
	err := r.collection.FindOne(ctx, filter).Decode(&chars)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("characteristics not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get characteristics: %w", err)
	}

	return &chars, nil
}

func (r *CharacteristicsRepository) Delete(ctx context.Context, listingID uuid.UUID) error {
	filter := bson.M{"listing_id": listingID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete characteristics: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("characteristics not found")
	}

	return nil
}

func (r *CharacteristicsRepository) AddTag(ctx context.Context, listingID uuid.UUID, tag string) error {
	filter := bson.M{"listing_id": listingID}
	update := bson.M{
		"$addToSet": bson.M{"tags": tag},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

func (r *CharacteristicsRepository) RemoveTag(ctx context.Context, listingID uuid.UUID, tag string) error {
	filter := bson.M{"listing_id": listingID}
	update := bson.M{
		"$pull": bson.M{"tags": tag},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}
