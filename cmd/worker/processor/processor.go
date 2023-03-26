package processor

import (
	"context"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
	transactionsUseCase "github.com/Inspirate789/SOK-golang-test-task/internal/transactions/usecase"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/topics"
	"time"
)

type KafkaProcessorConfig struct {
	BrokerUrls      []string
	ClientID        string
	NewConsumerFunc func(cfg *consumer.KafkaConsumerConfig) consumer.TransactionConsumer
	DetectionTime   time.Duration
	UseCase         transactionsUseCase.UseCase
	Logger          *zerolog.Logger
}

type KafkaProcessor struct {
	client          *kafka.Client
	brokerUrls      []string
	clientID        string
	consumers       map[string]consumer.TransactionConsumer
	newConsumerFunc func(cfg *consumer.KafkaConsumerConfig) consumer.TransactionConsumer
	detectionTime   time.Duration
	useCase         transactionsUseCase.UseCase
	logger          *zerolog.Logger
}

func NewKafkaProcessor(cfg *KafkaProcessorConfig) *KafkaProcessor {
	return &KafkaProcessor{
		client: &kafka.Client{
			Addr:      kafka.TCP(cfg.BrokerUrls...),
			Timeout:   0,
			Transport: &kafka.Transport{ClientID: cfg.ClientID},
		},
		brokerUrls:      cfg.BrokerUrls,
		clientID:        cfg.ClientID,
		consumers:       make(map[string]consumer.TransactionConsumer),
		newConsumerFunc: cfg.NewConsumerFunc,
		detectionTime:   cfg.DetectionTime,
		useCase:         cfg.UseCase,
		logger:          cfg.Logger,
	}
}

func (kp *KafkaProcessor) getTopicsList(ctx context.Context) ([]string, error) {
	kafkaTopics, err := topics.List(ctx, kp.client)
	if err != nil {
		return nil, err
	}

	topicNames := make([]string, 0, len(kafkaTopics))
	for _, topic := range kafkaTopics {
		topicNames = append(topicNames, topic.Name)
	}

	return topicNames, nil
}

func (kp *KafkaProcessor) topicsDetectionRoutine(ctx context.Context, newTopics chan string) error {
	ticker := time.NewTicker(kp.detectionTime)
	for {
		select {
		case <-ticker.C:
			topicNames, err := kp.getTopicsList(ctx)
			if err != nil {
				ticker.Stop()
				ctx.Done()
				return err
			}
			for _, topic := range topicNames {
				if _, inMap := kp.consumers[topic]; !inMap {
					newTopics <- topic
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func (kp *KafkaProcessor) messageDetectionRoutine(ctx context.Context, consumer consumer.TransactionConsumer, transactions chan consumer.TransactionDTO) error {
	for {
		transaction, err := consumer.Consume(ctx)
		if err != nil {
			transactions <- transaction
		} else {
			ctx.Done()
			return err
		}
	}
}

func (kp *KafkaProcessor) consumerRoutine(ctx context.Context, cons consumer.TransactionConsumer) error {
	transactions := make(chan consumer.TransactionDTO)
	var detectionErr error
	go func() {
		detectionErr = kp.messageDetectionRoutine(ctx, cons, transactions)
	}()

	var processErr error
OUTER:
	for {
		select {
		case transaction := <-transactions:
			processErr = kp.useCase.PerformTransaction(transaction)
			if processErr != nil {
				ctx.Done()
				break OUTER
			}
		case <-ctx.Done():
			break OUTER
		}
	}

	if processErr != nil {
		kp.logger.Error().Err(ErrFinishWrap(processErr))
		return processErr
	} else if detectionErr != nil {
		kp.logger.Error().Err(ErrFinishWrap(detectionErr))
		return detectionErr
	}
	kp.logger.Debug().Msg("kafka processor finished successfully")

	return nil
}

func (kp *KafkaProcessor) ProcessQueue(ctx context.Context) error {
	newTopics := make(chan string)
	var detectionErr error
	go func() {
		detectionErr = kp.topicsDetectionRoutine(ctx, newTopics)
	}()

	var consumeErr error
OUTER:
	for {
		select {
		case newTopic := <-newTopics:
			kp.consumers[newTopic] = kp.newConsumerFunc(&consumer.KafkaConsumerConfig{
				BrokerUrls: kp.brokerUrls,
				ClientID:   kp.clientID,
				Topic:      newTopic,
			})
			go func() {
				consumeErr = kp.consumerRoutine(ctx, kp.consumers[newTopic])
			}()
		case <-ctx.Done():
			break OUTER
		}
	}

	if consumeErr != nil {
		kp.logger.Error().Err(ErrFinishWrap(consumeErr))
		return consumeErr
	} else if detectionErr != nil {
		kp.logger.Error().Err(ErrFinishWrap(detectionErr))
		return detectionErr
	}
	kp.logger.Debug().Msg("kafka processor finished successfully")

	return nil
}
