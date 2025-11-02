package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.TopicListingsEvents,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(cfg.RequiredAcks),
		WriteTimeout: cfg.WriteTimeout,
		MaxAttempts:  cfg.RetryMax,
	}

	log.Printf("Kafka producer initialized for topic: %s", cfg.TopicListingsEvents)

	return &Producer{
		writer: writer,
		topic:  cfg.TopicListingsEvents,
	}, nil
}

func (p *Producer) PublishEvent(ctx context.Context, event ListingEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(event.ListingID),
		Value: eventBytes,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	log.Printf("Event published: %s for listing: %s", event.EventType, event.ListingID)
	return nil
}

func (p *Producer) PublishListingCreated(ctx context.Context, userID, listingID, title string) error {
	event := ListingEvent{
		EventID:   generateEventID(),
		EventType: EventListingCreated,
		Timestamp: time.Now(),
		UserID:    userID,
		ListingID: listingID,
		Data: map[string]interface{}{
			"title":  title,
			"action": "Ваше объявление успешно создано",
		},
	}
	return p.PublishEvent(ctx, event)
}

func (p *Producer) PublishListingDeleted(ctx context.Context, userID, listingID, title string) error {
	event := ListingEvent{
		EventID:   generateEventID(),
		EventType: EventListingDeleted,
		Timestamp: time.Now(),
		UserID:    userID,
		ListingID: listingID,
		Data: map[string]interface{}{
			"title":  title,
			"action": "Ваше объявление удалено",
		},
	}
	return p.PublishEvent(ctx, event)
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("failed to close kafka writer: %w", err)
	}
	log.Println("Kafka producer closed successfully")
	return nil
}

func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
