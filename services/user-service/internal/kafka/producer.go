package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/config"
	"log/slog"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       *slog.Logger
}

func NewProducer(cfg config.KafkaConfig, logger *slog.Logger) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = cfg.RequiredAcks
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = cfg.CompressionType

	config.Net.DialTimeout = cfg.DialTimeout
	config.Net.WriteTimeout = cfg.WriteTimeout
	config.Net.ReadTimeout = cfg.ReadTimeout

	config.Version = sarama.V3_5_0_0

	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync producer: %w", err)
	}

	return &Producer{
		syncProducer: syncProducer,
		topic:        cfg.TopicUserEvents,
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
