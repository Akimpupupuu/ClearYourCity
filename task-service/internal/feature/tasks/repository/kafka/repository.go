package tasks_kafka

import (
	"context"
	"fmt"
	"strconv"

	core_kafka "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/kafka"
)

const (
	taskNotificationTopic = "task-created"
)

type tasksKafkaRepository struct {
	producer *core_kafka.Producer
}

func NewTasksKafkaRepository(producer *core_kafka.Producer) *tasksKafkaRepository {
	return &tasksKafkaRepository{
		producer: producer,
	}
}

func (r *tasksKafkaRepository) PublishTaskCreated(ctx context.Context, taskID int, message []byte) error {
	key := []byte(strconv.Itoa(taskID))
	if err := r.producer.Produce(ctx, taskNotificationTopic, key, message); err != nil {
		return fmt.Errorf("publish message to kafka: %w", err)
	}

	return nil
}
