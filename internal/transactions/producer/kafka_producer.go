package producer

import (
	"context"
	"encoding/json"
	"github.com/Inspirate789/SOK-golang-test-task/internal/models"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"time"
)

type KafkaProducerConfig struct {
	BrokerUrls []string
	ClientID   string
}

type KafkaTransactionProducer struct {
	writer     *kafka.Writer
	topicNames map[string]bool
	brokerUrls []string
}

func NewKafkaTransactionProducer(cfg *KafkaProducerConfig) TransactionProducer {
	return &KafkaTransactionProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.BrokerUrls...),
			Balancer:     &kafka.LeastBytes{},
			Transport:    &kafka.Transport{ClientID: cfg.ClientID},
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
			Compression:  compress.Snappy,
		},
		topicNames: make(map[string]bool),
		brokerUrls: cfg.BrokerUrls,
	}
}

func (ktp *KafkaTransactionProducer) createTopic(ctx context.Context, topicName string) (*kafka.Conn, error) {
	return kafka.DialLeader(ctx, "tcp", ktp.brokerUrls[0], topicName, 0)
}

func (ktp *KafkaTransactionProducer) Produce(ctx context.Context, queueName string, transaction consumer.TransactionDTO) error {
	if _, inMap := ktp.topicNames[queueName]; !inMap {
		_, err := ktp.createTopic(ctx, queueName)
		if err != nil {
			return models.ErrCreateTopicWrap(err)
		}
	}

	msg, err := json.Marshal(transaction)
	if err != nil {
		return models.ErrEncodingMsgWrap(err)
	}

	message := kafka.Message{
		Key:   []byte{byte(transaction.UserId)},
		Value: msg,
		Topic: queueName,
		Time:  time.Now(),
	}

	err = ktp.writer.WriteMessages(ctx, message)
	if err != nil {
		return models.ErrSendingMsgWrap(err)
	}
	// log.Debug().Msg(fmt.Sprintf("Send message: %v", string(msg)))
	return nil
}

func (ktp *KafkaTransactionProducer) Close() error {
	return ktp.writer.Close()
}
