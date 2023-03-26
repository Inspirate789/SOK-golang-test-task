package consumer

import (
	"context"
	"encoding/json"
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaConsumerConfig struct {
	BrokerUrls []string
	ClientID   string
	Topic      string
}

type KafkaTransactionConsumer struct {
	reader *kafka.Reader
}

func NewKafkaTransactionConsumer(cfg *KafkaConsumerConfig) TransactionConsumer {
	config := kafka.ReaderConfig{
		Brokers:         cfg.BrokerUrls,
		GroupID:         cfg.ClientID,
		Topic:           cfg.Topic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	}

	return &KafkaTransactionConsumer{
		reader: kafka.NewReader(config),
	}
}

func (ktc *KafkaTransactionConsumer) Consume(ctx context.Context) (TransactionDTO, error) {
	msg, err := ktc.reader.ReadMessage(context.Background())
	if err != nil {
		return TransactionDTO{}, models.ErrReceivingMsgWrap(err)
	}

	var transaction TransactionDTO
	err = json.Unmarshal(msg.Value, &transaction)
	if err != nil {
		return TransactionDTO{}, models.ErrDecodingMsgWrap(err)
	}

	// TODO: add log record about receiving message

	return transaction, nil
}

func (ktc *KafkaTransactionConsumer) Close() error {
	return ktc.reader.Close()
}
