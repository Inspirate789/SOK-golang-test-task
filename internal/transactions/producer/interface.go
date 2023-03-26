package producer

import (
	"context"
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
)

type TransactionProducer interface {
	Produce(ctx context.Context, queueName string, transaction consumer.TransactionDTO) error
	Close() error
}
