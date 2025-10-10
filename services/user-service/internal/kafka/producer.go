package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"time"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       *slog.Logger
}

func NewProducer(brokers []string, topic string, logger *slog.Logger) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy

	config.Net.DialTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second

	config.Version = sarama.V3_5_0_0

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync producer: %w", err)
	}

	return &Producer{
		syncProducer: syncProducer,
		topic:        topic,
		logger:       logger,
	}, nil
}

func (p *Producer) SendEmailVerification(event *EmailVerificationEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal email verification event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.UserID),
		Value: sarama.ByteEncoder(eventJSON),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Failed to send message to Kafka",
			"error", err,
			"topic", p.topic,
			"user_id", event.UserID,
		)
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	p.logger.Debug("Message sent to Kafka",
		"topic", p.topic,
		"partition", partition,
		"offset", offset,
		"user_id", event.UserID,
		"event_type", event.EventType,
	)

	return nil
}

func (p *Producer) Close() error {
	if p.syncProducer != nil {
		return p.syncProducer.Close()
	}
	return nil
}
