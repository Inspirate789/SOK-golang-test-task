package repository

import (
	"github.com/Inspirate789/SOK-golang-test-task/internal/transactions/consumer"
)

type Repository interface {
	PerformTransaction(transaction consumer.TransactionDTO) error
}
