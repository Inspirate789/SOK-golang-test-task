package usecase

import (
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
)

type UseCase interface {
	PerformTransaction(transaction consumer.TransactionDTO) error
}
