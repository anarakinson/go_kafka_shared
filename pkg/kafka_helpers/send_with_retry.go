package kafka_helpers

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

// отправка с ретраем
func SendWithRetry(producer sarama.SyncProducer, message *sarama.ProducerMessage, maxRetries int) (int32, int64, error) {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		partition, offset, err := producer.SendMessage(message)
		if err == nil {
			// Успешная отправка
			return partition, offset, nil
		}

		// проверяем ошибку от кафки
		if kafkaError, ok := err.(sarama.ProducerError); ok {
			if kafkaError.Err == sarama.ErrLeaderNotAvailable ||
				kafkaError.Err == sarama.ErrNotEnoughReplicas ||
				kafkaError.Err == sarama.ErrNotEnoughReplicasAfterAppend {
				// временные ошибки - повторяем
				logger.Log.Error("Error sending message", zap.Int("Attempt", attempt), zap.Error(err))
				time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
				continue
			}
		}

		// Постоянная ошибка - не повторяем
		break
	}
	return 0, 0, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}
