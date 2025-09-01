package kafka_helpers

import (
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
)

var (
	ErrCanceled    = errors.New("context canceled")
	ErrCreateAdmin = errors.New("error creating kafka admin")
	ErrListTopics  = errors.New("error listing kafka topics")
	ErrCreateTopic = errors.New("error creating topic")
	ErrTopicExists = errors.New("topic already exists")
)

// создание топика в кафке,
func CreateTopic(brokers []string, topicName string) error {
	admin, err := sarama.NewClusterAdmin(
		brokers, // адреса узлов
		sarama.NewConfig(),
	)
	if err != nil {
		slog.Error("Error creating sarama admin client", "error", err)
		return ErrCreateAdmin
	}
	defer admin.Close()

	// получаем список существующих топиков
	topics, err := admin.ListTopics()
	if err != nil {
		slog.Error("Error listing topics", "error", err)
		return ErrListTopics
	}

	// проверяем существование топика
	_, exists := topics[topicName]
	if exists {
		slog.Info("Topic already exists", "topic name", topicName)
		return ErrTopicExists
	}

	// если не существует
	// создаем топик с 2 партициями и фактором репликации 1
	err = admin.CreateTopic(
		topicName,
		&sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		},
		false,
	)
	if err != nil {
		slog.Error("Error creating topic", "topic name", topicName, "error", err)
		return ErrCreateTopic
	}

	return nil
}
