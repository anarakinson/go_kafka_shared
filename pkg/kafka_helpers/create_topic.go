package kafka_helpers

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
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
		logger.Log.Error("Error creating sarama admin client", zap.Error(err))
		return ErrCreateAdmin
	}
	defer admin.Close()

	// получаем список существующих топиков
	topics, err := admin.ListTopics()
	if err != nil {
		logger.Log.Error("Error listing topics", zap.Error(err))
		return ErrListTopics
	}

	// проверяем существование топика
	_, exists := topics[topicName]
	if exists {
		logger.Log.Info("Topic already exists", zap.String("topic name", topicName))
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
		logger.Log.Error("Error creating topic", zap.String("topic name", topicName), zap.Error(err))
		return ErrCreateTopic
	}

	return nil
}
