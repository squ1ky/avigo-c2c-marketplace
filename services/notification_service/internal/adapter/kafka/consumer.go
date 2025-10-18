package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/config"
)

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

type MessageHandler interface {
	Handle(ctx context.Context, message kafka.Message) error
}

func NewConsumer(cfg *config.KafkaConfig, handler MessageHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.TopicUserEvents,
		GroupID:        cfg.GroupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		Logger:         kafka.LoggerFunc(log.Printf),
		ErrorLogger:    kafka.LoggerFunc(log.Printf),
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Start Kafka consumer...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping Kafka consumer...")
			return c.Close()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return c.Close()
				}
				log.Printf("Erro fetching message from Kafka: %v", err)
				time.Sleep(time.Second)
				continue
			}

			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	log.Printf("Processing message from partition %d, offset %d", msg.Partition, msg.Offset)

	processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := c.handler.Handle(processCtx, msg); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}

	return nil
}

func (c *Consumer) Close() error {
	log.Println("Closing Kafka consumer...")
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close kafka reader: %w", err)
	}
	log.Println("Kafka consumer closed successfully")
	return nil
}
