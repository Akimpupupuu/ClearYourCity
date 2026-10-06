package core_kafka

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(ctx context.Context, config Config, logger core_logger.Logger) (*Producer, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided: %w", core_errors.ErrInvalidArgument)
	}
	if len(config.Topics) == 0 {
		return nil, fmt.Errorf("no topics provided: %w", core_errors.ErrInvalidArgument)
	}

	var conn *kafka.Conn
	var err error
	for _, broker := range config.Brokers {
		conn, err = kafka.DialContext(ctx, "tcp", broker)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("create a dial connection to any kafka broker: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("close connection", core_logger.Err(err))
		}
	}()

	for _, topic := range config.Topics {
		if topic == "" {
			continue
		}

		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     config.NumPartitions,
			ReplicationFactor: config.ReplicationFactor,
		})
		if err != nil && !errors.Is(err, kafka.TopicAlreadyExists) {
			return nil, fmt.Errorf("create topic: %s: %w", topic, err)
		}
	}

	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(config.Brokers...),
			Async: false,
		},
	}, nil
}

func (p *Producer) Produce(ctx context.Context, topic string, key []byte, value []byte) error {
	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}); err != nil {
		return fmt.Errorf("write kafka to topic: %s: %w", topic, err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("close kafka writer: %w", err)
	}

	return nil
}
